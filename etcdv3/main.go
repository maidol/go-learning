package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"

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
	fmt.Printf("%v\n", config)
	cli, err := etcd.New(config)

	rch := cli.Watch(context.Background(), "sample_key")
	for wresp := range rch {
		fmt.Println("len(wresp.Events)", len(wresp.Events))
		fmt.Printf("wresp %+v\n", wresp)
		for _, ev := range wresp.Events {
			fmt.Printf("ev %+v\n", ev)
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if string(ev.Kv.Key) != "sample_key" || ev.Type != mvccpb.PUT {
				// continue
			}
		}
	}

	put(cli, "sample_key", "sample_value")

	defer cli.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
