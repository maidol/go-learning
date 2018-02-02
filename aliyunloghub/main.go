package main

import (
	"aliyun-log-go-sdk/example/util"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
)

func main() {
	logproject := &sls.LogProject{}

	flag.StringVar(&logproject.Name, "projectname", "epaper", "loghub project name")
	flag.StringVar(&logproject.Endpoint, "endpoint", "cn-beijing.log.aliyuncs.com", "loghub endpoint")
	flag.StringVar(&logproject.AccessKeyID, "accesskeyid", "", "loghub AccessKeyID")
	flag.StringVar(&logproject.AccessKeySecret, "accesskeysecret", "", "loghub AccessKeySecret")

	flag.Parse()

	util.Project = logproject
	fmt.Println("loghub sample begin")
	begin_time := uint32(time.Now().Unix())
	rand.Seed(int64(begin_time))
	logstore_name := "gateway"
	var retry_times int
	var logstore *sls.LogStore
	var err error
	for retry_times = 0; ; retry_times++ {
		if retry_times > 5 {
			return
		}
		logstore, err = util.Project.GetLogStore(logstore_name)
		if err != nil {
			fmt.Printf("GetLogStore fail, retry:%d, err:%v\n", retry_times, err)
			if strings.Contains(err.Error(), sls.PROJECT_NOT_EXIST) {
				return
			} else if strings.Contains(err.Error(), sls.LOGSTORE_NOT_EXIST) {
				err = util.Project.CreateLogStore(logstore_name, 1, 2)
				if err != nil {
					fmt.Printf("CreateLogStore fail, err: ", err.Error())
				} else {
					fmt.Println("CreateLogStore success")
				}
			}
		} else {
			fmt.Printf("GetLogStore success, retry:%d, name: %s, ttl: %d, shardCount: %d, createTime: %d, lastModifyTime: %d\n", retry_times, logstore.Name, logstore.TTL, logstore.ShardCount, logstore.CreateTime, logstore.LastModifyTime)
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	// put logs to logstore
	for loggroupIdx := 0; loggroupIdx < 1; loggroupIdx++ {
		logs := []*sls.Log{}
		for logIdx := 0; logIdx < 1; logIdx++ {
			content := []*sls.LogContent{}
			for colIdx := 0; colIdx < 5; colIdx++ {
				content = append(content, &sls.LogContent{
					Key:   proto.String(fmt.Sprintf("col_%d", colIdx)),
					Value: proto.String(fmt.Sprintf("loggroup idx: %d, log idx: %d, col idx: %d, value: %d", loggroupIdx, logIdx, colIdx, rand.Intn(10000000))),
				})
			}
			log := &sls.Log{
				Time:     proto.Uint32(uint32(time.Now().Unix())),
				Contents: content,
			}
			logs = append(logs, log)
		}
		loggroup := &sls.LogGroup{
			Topic:  proto.String("gateway"),
			Source: proto.String("10.230.201.117"),
			Logs:   logs,
		}
		// PostLogStoreLogs API Ref: https://intl.aliyun.com/help/doc-detail/29026.htm
		for retry_times = 0; retry_times < 10; retry_times++ {
			err := logstore.PutLogs(loggroup)
			if err == nil {
				fmt.Printf("PutLogs success, retry: %d\n", retry_times)
				break
			} else {
				fmt.Printf("PutLogs fail, retry: %d, err: %s\n", retry_times, err)
				//handle exception here, you can add retryable erorrCode, set appropriate put_retry
				if strings.Contains(err.Error(), sls.WRITE_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_WRITE_QUOTA_EXCEED) {
					//mayby you should split shard
					time.Sleep(1000 * time.Millisecond)
				} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
		time.Sleep(200 * time.Millisecond)
	}

	// // pull logs from logstore
	// var shards []int
	// for retry_times = 0; ; retry_times++ {
	// 	if retry_times > 5 {
	// 		return
	// 	}
	// 	shards, err = logstore.ListShards()
	// 	if err != nil {
	// 		fmt.Printf("ListShards fail, retry: %d, err:%v\n", retry_times, err)
	// 	} else {
	// 		fmt.Printf("ListShards success, shards: %v\n", shards)
	// 		break
	// 	}
	// 	time.Sleep(200 * time.Millisecond)
	// }
	// var begin_cursor string
	// var end_cursor string
	// var next_cursor string
	// var loggrouplist *sls.LogGroupList
	// for _, sh := range shards {
	// 	if sh == 0 {
	// 		// sample of pulllogs from begin
	// 		// GetCursor API Ref: https://intl.aliyun.com/help/doc-detail/29024.htm
	// 		for retry_times = 0; ; retry_times++ {
	// 			if retry_times > 5 {
	// 				return
	// 			}
	// 			begin_cursor, err = logstore.GetCursor(sh, "begin")
	// 			if err != nil {
	// 				fmt.Printf("GetCursor(begin) fail, retry: %d, err:%v\n", retry_times, err)
	// 			} else {
	// 				fmt.Printf("GetCursor(begin) success\n")
	// 				break
	// 			}
	// 			time.Sleep(200 * time.Millisecond)
	// 		}
	// 		for retry_times = 0; ; retry_times++ {
	// 			if retry_times > 5 {
	// 				return
	// 			}
	// 			end_cursor, err = logstore.GetCursor(sh, "end")
	// 			if err != nil {
	// 				fmt.Printf("GetCursor(end) fail, retry: %d, err:%v\n", retry_times, err)
	// 			} else {
	// 				fmt.Printf("GetCursor(end) success\n")
	// 				break
	// 			}
	// 			time.Sleep(200 * time.Millisecond)
	// 		}
	// 		// PullLogs API Ref: https://intl.aliyun.com/help/doc-detail/29025.htm
	// 		for retry_times = 0; ; retry_times++ {
	// 			if retry_times > 100 {
	// 				return
	// 			}
	// 			loggrouplist, next_cursor, err = logstore.PullLogs(sh, begin_cursor, end_cursor, 100)
	// 			if err == nil {
	// 				fmt.Printf("PullLogs success, retry: %d\n", retry_times)
	// 				break
	// 			} else {
	// 				fmt.Printf("PullLogs fail frome shard %v, retry: %d, err: %s\n", sh, retry_times, err)
	// 				//handle exception here, you can add retryable erorrCode, set appropriate put_retry
	// 				if strings.Contains(err.Error(), sls.READ_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_READ_QUOTA_EXCEED) {
	// 					//mayby you should split shard
	// 					time.Sleep(1000 * time.Millisecond)
	// 				} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
	// 					time.Sleep(200 * time.Millisecond)
	// 				}
	// 			}
	// 		}
	// 		fmt.Printf("shard: %d, begin_cursor: %s, end_cursor: %s, next_cursor: %s\n", sh, begin_cursor, end_cursor, next_cursor)
	// 		for _, loggroup := range loggrouplist.LogGroups {
	// 			for _, log := range loggroup.Logs {
	// 				for _, content := range log.Contents {
	// 					if content != nil {
	// 						fmt.Printf("key:%s, value:%s\n", content.GetKey(), content.GetValue())
	// 					}
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		// sample of pulllogs from setted time
	// 		for retry_times = 0; ; retry_times++ {
	// 			if retry_times > 5 {
	// 				return
	// 			}
	// 			begin_cursor, err = logstore.GetCursor(sh, strconv.Itoa(int(begin_time)+2))
	// 			if err != nil {
	// 				fmt.Printf("GetCursor fail, retry: %d, err:%v\n", retry_times, err)
	// 			} else {
	// 				fmt.Printf("GetCursor success\n")
	// 				fmt.Printf("begin_cursor: %s\n", begin_cursor)
	// 				break
	// 			}
	// 			time.Sleep(200 * time.Millisecond)
	// 		}
	// 		for {
	// 			for retry_times = 0; ; retry_times++ {
	// 				if retry_times > 2 {
	// 					return
	// 				}
	// 				loggrouplist, next_cursor, err = logstore.PullLogs(sh, begin_cursor, "", 2)
	// 				if err == nil {
	// 					fmt.Printf("PullLogs success, retry: %d\n", retry_times)
	// 					break
	// 				} else {
	// 					fmt.Printf("PullLogs fail frome shard %v, retry: %d, err: %s\n", sh, retry_times, err)
	// 					//handle exception here, you can add retryable erorrCode, set appropriate put_retry
	// 					if strings.Contains(err.Error(), sls.READ_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_READ_QUOTA_EXCEED) {
	// 						//mayby you should split shard
	// 						time.Sleep(1000 * time.Millisecond)
	// 					} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
	// 						time.Sleep(200 * time.Millisecond)
	// 					}
	// 				}
	// 			}
	// 			fmt.Printf("shard: %d, begin_cursor: %s, next_cursor: %s, len(loggrouplist.LogGroups): %d\n", sh, begin_cursor, next_cursor, len(loggrouplist.LogGroups))
	// 			if len(loggrouplist.LogGroups) == 0 {
	// 				// means no more data in this shard, you can break out or sleep to wait new data
	// 				break
	// 			} else {
	// 				for _, loggroup := range loggrouplist.LogGroups {
	// 					for _, log := range loggroup.Logs {
	// 						for _, content := range log.Contents {
	// 							fmt.Printf("key:%s, value:%s\n", content.GetKey(), content.GetValue())
	// 						}
	// 					}
	// 				}
	// 				begin_cursor = next_cursor
	// 			}
	// 		}
	// 	}
	// }
	fmt.Println("loghub sample end")
}
