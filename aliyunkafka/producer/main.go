package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"go-learning/aliyunkafka"

	"github.com/Shopify/sarama"
)

var cfg *configs.MqConfig
var producer sarama.SyncProducer

func init() {

	fmt.Print("init kafka producer\n")

	var err error

	cfg = &configs.MqConfig{}
	configs.LoadJsonConfig(cfg, "mq.json")

	flag.StringVar(&cfg.Ak, "ak", "", "access key")
	flag.StringVar(&cfg.Password, "password", "", "password")
	flag.Parse()

	fmt.Printf("load config: %v\n", cfg)

	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = cfg.Ak
	mqConfig.Net.SASL.Password = cfg.Password
	mqConfig.Net.SASL.Handshake = true

	certBytes, err := ioutil.ReadFile(configs.GetFullPath(cfg.CertFile))
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka producer failed to parse root certificate")
	}

	mqConfig.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Successes = true

	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *cfg, err)
		fmt.Println(msg)
		panic(msg)
	}

	producer, err = sarama.NewSyncProducer(cfg.Servers, mqConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		fmt.Println(msg)
		panic(msg)
	}

}

func produce(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		msg := fmt.Sprintf("Kafka send message error. topic: %v. key: %v. content: %v", topic, key, content)
		fmt.Println(msg)
		return err
	}

	return nil
}

func main() {
	key := time.Now().UnixNano()
	produce(cfg.Topics[0], strconv.FormatInt(key, 10), "this is a kafka message!!!!!")
	fmt.Printf("Send OK key:%s value:%s\n", strconv.FormatInt(key, 10), "this is a kafka message!!!!!")
}
