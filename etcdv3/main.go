package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"

	etcd "github.com/coreos/etcd/clientv3"
)

func put(cli *etcd.Client, key, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceld by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)

		}
	}
}

func main() {
	config := etcd.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	// fmt.Println(config)
	cli, err := etcd.New(config)
	fmt.Printf("%v", config)

	put(cli, "sample_key", "sample_value")

	defer cli.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
