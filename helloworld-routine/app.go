package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched() // 让CPU把时间片让给别人,下次某个时候继续恢复执行该goroutine
		fmt.Println(s)
	}
}

func main() {
	fmt.Println(runtime.NumGoroutine())
	go say("world0") // 开一个新的Goroutines执行
	go say("world1") // 开一个新的Goroutines执行
	go say("world2") // 开一个新的Goroutines执行
	// go say("world3") // 开一个新的Goroutines执行
	say("hello") // 当前Goroutines执行
	fmt.Println(runtime.NumGoroutine())
}
