package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"os/signal"

	"go-learning/aliyunkafka"

	"github.com/Shopify/sarama"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/bsm/sarama-cluster"
	"github.com/gogo/protobuf/proto"
)

var cfg *configs.MqConfig
var consumer *cluster.Consumer
var sig chan os.Signal

var loghub *Loghub

func init() {
	fmt.Println("init kafka consumer")

	var err error

	// loghub
	lhcfg := &Config{}
	flag.StringVar(&lhcfg.LogProject.Name, "projectname", "epaper", "loghub project name")
	flag.StringVar(&lhcfg.LogProject.Endpoint, "endpoint", "cn-beijing.log.aliyuncs.com", "loghub endpoint")
	flag.StringVar(&lhcfg.LogProject.AccessKeyID, "accesskeyid", "", "loghub AccessKeyID")
	flag.StringVar(&lhcfg.LogProject.AccessKeySecret, "accesskeysecret", "", "loghub AccessKeySecret")

	cfg := &configs.MqConfig{}
	configs.LoadJsonConfig(cfg, "mq.json")
	flag.StringVar(&cfg.Ak, "ak", "", "access key")
	flag.StringVar(&cfg.Password, "password", "", "password")
	flag.Parse()

	fmt.Printf("load config: %v\n", cfg)

	clusterCfg := cluster.NewConfig()

	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = cfg.Ak
	clusterCfg.Net.SASL.Password = cfg.Password
	clusterCfg.Net.SASL.Handshake = true

	certBytes, err := ioutil.ReadFile(configs.GetFullPath(cfg.CertFile))
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka consumer failed to parse root certificate")
	}

	clusterCfg.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	clusterCfg.Net.TLS.Enable = true
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	clusterCfg.Group.Return.Notifications = true

	clusterCfg.ChannelBufferSize = 1024

	clusterCfg.Version = sarama.V0_10_0_0
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		fmt.Println(msg)
		panic(msg)
	}

	consumer, err = cluster.NewConsumer(cfg.Servers, cfg.ConsumerId, cfg.Topics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		fmt.Println(msg)
		panic(msg)
	}

	sig = make(chan os.Signal, 1)

	// loghub
	loghub = NewLoghub(lhcfg, consumer)
	loghub.Run()
}

func Start() {
	go consume()
}

func consume() {
	for {
		select {
		case msg, more := <-consumer.Messages():
			if more {
				fmt.Printf("kafka consumer msg: (topic:%s) (partition:%d) (offset:%d) (%s): (%s)\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				loghub.Input() <- msg
				// consumer.MarkOffset(msg, "completed") // mark message as processed
				// fmt.Println("kafka consumer HighWaterMarks", consumer.HighWaterMarks())
			}
		case err, more := <-consumer.Errors():
			if more {
				fmt.Printf("Kafka consumer error: %v\n", err.Error())
			}
		case ntf, more := <-consumer.Notifications():
			if more {
				fmt.Printf("Kafka consumer rebalance: %v\n", ntf)
			}
		case <-sig:
			fmt.Errorf("Stop consumer server...")
			consumer.Close()
			return
		}
	}

}

func Stop(s os.Signal) {
	fmt.Println("Recived kafka consumer stop signal...")
	sig <- s
	fmt.Println("kafka consumer stopped!!!")
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	Start()

	select {
	case s := <-signals:
		Stop(s)
	}

}

type Config struct {
	Name       string
	LogProject struct {
		Name            string
		Endpoint        string
		AccessKeyID     string
		AccessKeySecret string
	}
	MessageChannelBufferSize int
	LogsBufferSize           int
	Topics                   []string
	LogsBufferSize4Logstore  int
	Logstores                []string
}

type Loghub struct {
	Name string
	*Config
	consumer                 *cluster.Consumer
	logproject               *sls.LogProject
	logstores                map[string]*sls.LogStore
	messages                 chan *sarama.ConsumerMessage
	messageChannelBufferSize int
	m                        sync.RWMutex
	stop                     chan int

	mlogstoreLogsBuffer     sync.RWMutex
	logsBuffer4Logstore     map[string](chan *topicLog) // by logstore
	logsBufferSize4Logstore int

	mlogsBuffer sync.RWMutex
	// logsBuffer     map[string](chan *topicLog) // by topic and logstore
	logsBuffer     map[string]map[string]chan *topicLog // by topic and logstore
	logsBufferSize int
	topics         []string
}

type topicLog struct {
	topic string
	log   *sls.Log
	cmsg  *sarama.ConsumerMessage
}

func NewLoghub(cfg *Config, consumer *cluster.Consumer) *Loghub {
	logproject := &sls.LogProject{
		Name:            cfg.LogProject.Name,
		Endpoint:        cfg.LogProject.Endpoint,
		AccessKeyID:     cfg.LogProject.AccessKeyID,
		AccessKeySecret: cfg.LogProject.AccessKeySecret,
	}

	lbls := map[string](chan *topicLog){}

	lbr := map[string]map[string](chan *topicLog){}
	// lbr := map[string](chan *topicLog){}
	// for _, tp := range cfg.Topics {
	// 	lbr[tp] = make(chan *topicLog, cfg.LogsBufferSize)
	// }

	lh := &Loghub{
		Name:                     cfg.Name,
		Config:                   cfg,
		consumer:                 consumer,
		logproject:               logproject,
		messages:                 make(chan *sarama.ConsumerMessage, cfg.MessageChannelBufferSize),
		messageChannelBufferSize: cfg.MessageChannelBufferSize,
		stop: make(chan int),

		logsBuffer4Logstore:     lbls,
		logsBufferSize4Logstore: cfg.LogsBufferSize4Logstore,
		logstores:               map[string]*sls.LogStore{},

		logsBuffer:     lbr,
		logsBufferSize: cfg.LogsBufferSize,
		topics:         cfg.Topics,
	}

	return lh
}

func (l *Loghub) Run() {
	// lss, err := l.logproject.ListLogStore()
	// if err != nil {
	// 	panic(err)
	// }
	// 开启日志库
	for _, lsn := range l.Logstores {
		_, err := l.getLogstore(lsn)
		if err != nil {
			fmt.Printf("Loghub Start failed (logstoreName=%s). err: %v\n", lsn, err)
			panic(err)
		}
		// 分配到topic
		go l.dispatchToTopic(lsn)
		for _, tp := range l.topics {
			go l.processTopicMsg(lsn, tp)
		}
	}

	// 分配消息
	go l.dispatch()
}

func (l *Loghub) Input() chan<- *sarama.ConsumerMessage {
	return l.messages
}

// dispatchToTopic
func (l *Loghub) dispatchToTopic(logstoreName string) {
	// TODO: 处理消息, 分配到不同的topic
	channelBuffer := l.getLogstoreLogsBufferChannel(logstoreName)
	for {
		select {
		case log := <-channelBuffer:
			l.mlogsBuffer.Lock()
			logsCB, ok := l.logsBuffer[logstoreName]
			if !ok || logsCB == nil {
				logsCB = map[string]chan *topicLog{}
				l.logsBuffer[logstoreName] = logsCB
			}
			logsCBTopic, ok := logsCB[log.topic]
			if !ok || logsCBTopic == nil {
				logsCBTopic = make(chan *topicLog, l.logsBufferSize)
				logsCB[log.topic] = logsCBTopic
			}
			l.mlogsBuffer.Unlock()
			logsCBTopic <- log
		}
	}
}

// dispatch
// 分配消息
func (l *Loghub) dispatch() error {
	// TODO: logproject, logstore, topic
	// 指定logproject和logstore进行分配
	for {
		select {
		case msg := <-l.messages:
			data, err := unserialize(msg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			logprojectName, ok1 := data["logproject"]
			if !ok1 || logprojectName == "" {
				fmt.Println("loghub.dispatch: data[\"logproject\"] was empty")
				continue
			}
			logstoreName, ok2 := data["logstore"]
			if !ok2 || logstoreName == "" {
				fmt.Println("loghub.dispatch: data[\"logstore\"] was empty")
				continue
			}
			topic, ok3 := data["topic"]
			if !ok3 || topic == "" {
				fmt.Println("loghub.dispatch: data[\"topic\"] was empty")
				continue
			}
			log, err := generateLog(data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// logstore, err := l.getLogstore(logstoreName)
			// if err != nil {
			// 	fmt.Println(err)
			// 	continue
			// }
			log.cmsg = msg
			lblsc := l.getLogstoreLogsBufferChannel(logstoreName)
			select {
			// TODO: 考虑优化处理, lblsc如果满了的情况
			case lblsc <- log:
			}
		}
	}
}

func (l *Loghub) Stop() {
	l.stop <- 0
}

func (l *Loghub) processTopicMsg(topic string, logstoreName string) error {
	cb := l.logsBuffer[logstoreName][topic]
	for {
		select {
		case log := <-cb:
			loggroup := generateLoggroupByTopicLog(log, "")
			logstore, err := l.getLogstore(logstoreName)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = putLogs(logstore, loggroup)
			if err != nil {
				fmt.Println(err)
				continue
			}
			l.consumer.MarkOffset(log.cmsg, "loghub.processTopicMsg")
		}
	}
}

func (l *Loghub) getLogstore(logstoreName string) (*sls.LogStore, error) {
	var logstore *sls.LogStore
	l.m.RLock()
	logstore = l.logstores[logstoreName]
	l.m.RUnlock()
	if logstore != nil {
		return logstore, nil
	}
	var err error
	for retry_times := 0; ; retry_times++ {
		if retry_times > 5 {
			return nil, errors.New("GetLogStore retry_times > 5")
		}
		logstore, err = l.logproject.GetLogStore(logstoreName)
		if err != nil {
			fmt.Printf("GetLogStore fail, retry:%d, err:%v\n", retry_times, err)
			if strings.Contains(err.Error(), sls.PROJECT_NOT_EXIST) {
				return nil, err
			} else if strings.Contains(err.Error(), sls.LOGSTORE_NOT_EXIST) {
				err = l.logproject.CreateLogStore(logstoreName, 1, 2)
				if err != nil {
					fmt.Printf("CreateLogStore fail, err: ", err.Error())
					return nil, err
				} else {
					fmt.Println("CreateLogStore success")
					l.m.Lock()
					l.logstores[logstoreName] = logstore
					l.m.Unlock()
					return logstore, nil
				}
			}
		} else {
			fmt.Printf("GetLogStore success, retry:%d, name: %s, ttl: %d, shardCount: %d, createTime: %d, lastModifyTime: %d\n", retry_times, logstore.Name, logstore.TTL, logstore.ShardCount, logstore.CreateTime, logstore.LastModifyTime)
			l.m.Lock()
			l.logstores[logstoreName] = logstore
			l.m.Unlock()
			return logstore, nil
		}
	}
}

func (l *Loghub) getLogstoreLogsBufferChannel(logstoreName string) chan *topicLog {
	l.mlogstoreLogsBuffer.RLock()
	slslchan, ok := l.logsBuffer4Logstore[logstoreName]
	l.mlogstoreLogsBuffer.RUnlock()
	if !ok || slslchan == nil {
		slslchan = make(chan *topicLog, l.logsBufferSize4Logstore)
		l.mlogstoreLogsBuffer.Lock()
		l.logsBuffer4Logstore[logstoreName] = slslchan
		l.mlogstoreLogsBuffer.Unlock()
	}
	return slslchan
}

func generateLoggroupByTopicLog(tlog *topicLog, source string) *sls.LogGroup {
	logs := []*sls.Log{tlog.log}
	loggroup := &sls.LogGroup{
		Topic:  proto.String(tlog.topic),
		Source: proto.String(source),
		Logs:   logs,
	}
	return loggroup
}

func unserialize(msg *sarama.ConsumerMessage) (map[string]string, error) {
	var err error
	data := map[string]string{}
	err = json.Unmarshal(msg.Value, &data)
	if err != nil {
		fmt.Printf("[unserialize sarama.ConsumerMessage] json.Unmarshal err: %v\n", err)
		return nil, err
	}
	return data, nil
}

func generateLog(data map[string]string) (*topicLog, error) {
	contents := []*sls.LogContent{}
	for k, v := range data {
		contents = append(contents, &sls.LogContent{
			Key:   proto.String(k),
			Value: proto.String(v),
		})
	}
	t, e := time.Parse("2006-01-02T15:04:05+08:00", data["@timestamp"])
	if e != nil {
		t = time.Now()
	}
	log := &sls.Log{
		Time:     proto.Uint32(uint32(t.Unix())),
		Contents: contents,
	}
	tplog := &topicLog{
		topic: data["topic"],
		log:   log,
	}
	return tplog, nil
}

func putLogs(logstore *sls.LogStore, loggroup *sls.LogGroup) error {
	var retry_times int
	var err error
	// PostLogStoreLogs API Ref: https://intl.aliyun.com/help/doc-detail/29026.htm
	for retry_times = 0; retry_times < 10; retry_times++ {
		err = logstore.PutLogs(loggroup)
		if err == nil {
			fmt.Printf("PutLogs success, retry: %d\n", retry_times)
			return nil
		}
		fmt.Printf("PutLogs fail, retry: %d, err: %s\n", retry_times, err)
		//handle exception here, you can add retryable erorrCode, set appropriate put_retry
		if strings.Contains(err.Error(), sls.WRITE_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_WRITE_QUOTA_EXCEED) {
			//you should split shard
		} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
		}
		continue
	}
	return err
}

// func PutMsg(msg *sarama.ConsumerMessage) error {
// 	var err error
// 	data := map[string]string{}
// 	err = json.Unmarshal(msg.Value, &data)
// 	if err != nil {
// 		fmt.Printf("[PutMsg] json.Unmarshal err: ", err)
// 		return err
// 	}

// 	loggroup := generateLoggroupByTopic([]map[string]string{data}, msg.Topic, "")
// 	err = putLogs(logstore, loggroup)
// 	if err != nil {
// 		fmt.Printf("loghub putlogs failed. err:%v\n", err)
// 		return err
// 	}
// 	return nil
// }

// func PutMsgs(msgs []*sarama.ConsumerMessage) {
// 	fmt.Println("loghub sample begin")

// 	var err error

// 	// put logs to logstore
// 	datas := []map[string]string{
// 		{"topic": "gateway-error", "key": strconv.FormatInt(time.Now().UnixNano(), 10), "message": "hello world!"},
// 		{"topic": "gateway-error", "key": strconv.FormatInt(time.Now().UnixNano(), 10), "message": "hello world!"},
// 		{"topic": "gateway-error", "key": strconv.FormatInt(time.Now().UnixNano(), 10), "message": "hello world!"},
// 	}
// 	loggroup := generateLoggroupByTopic(datas, "gateway-error")
// 	err = putLogs(logstore, loggroup)
// 	if err != nil {
// 		fmt.Printf("loghub putlogs failed. err:%v\n", err)
// 	}

// 	fmt.Println("loghub sample end")
// }

// func generateLoggroupByTopic(datas []map[string]string, topic string, source string) *sls.LogGroup {
// 	logs := []*sls.Log{}
// 	for _, data := range datas {
// 		contents := []*sls.LogContent{}
// 		for k, v := range data {
// 			contents = append(contents, &sls.LogContent{
// 				Key:   proto.String(k),
// 				Value: proto.String(v),
// 			})
// 		}
// 		log := &sls.Log{
// 			Time:     proto.Uint32(uint32(time.Now().Unix())),
// 			Contents: contents,
// 		}
// 		logs = append(logs, log)
// 	}
// 	loggroup := &sls.LogGroup{
// 		Topic:  proto.String(topic),
// 		Source: proto.String(source),
// 		Logs:   logs,
// 	}
// 	return loggroup
// }

// func generateLog(msg *sarama.ConsumerMessage) (*sls.Log, error) {
// 	var err error
// 	data := map[string]string{}
// 	err = json.Unmarshal(msg.Value, &data)
// 	if err != nil {
// 		fmt.Printf("[generateLog] json.Unmarshal err: %v\n", err)
// 		return nil, err
// 	}

// 	contents := []*sls.LogContent{}
// 	for k, v := range data {
// 		contents = append(contents, &sls.LogContent{
// 			Key:   proto.String(k),
// 			Value: proto.String(v),
// 		})
// 	}
// 	t, e := time.Parse("2006-01-02T15:04:05+08:00", data["@timestamp"])
// 	if e != nil {
// 		t = time.Now()
// 	}
// 	log := &sls.Log{
// 		Time:     proto.Uint32(uint32(t.Unix())),
// 		Contents: contents,
// 	}
// 	return log, nil
// }
