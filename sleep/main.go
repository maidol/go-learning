package main

import (
	"log"
	"runtime"
	"time"
)

// sleep 可以使当前goroutine暂时退出执行
func main() {
	runtime.GOMAXPROCS(1) // 限定1cpu(单核单线程)
	go func() {
		time.Sleep(1)
		log.Println("goroutine01")
	}()
	go func() {
		// (单核单线程)for{...}内部是空操作, 6s后主程序不会退出
		log.Println("goroutine02")
		for { // 不调用sleep或空操作, 空的 for{} 会完全阻塞当前goroutine所在的执行线程;貌似不遵循所说的时间片调度
			// time.Sleep(0) // 退出执行
			test()
		}
	}()
	go func() {
		time.Sleep(1)
		log.Println("goroutine03")
	}()
	go func() {
		log.Println("goroutine04")
	}()

	time.Sleep(6 * time.Second)
	log.Println("exit")
}

func test() {}
