package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	t := make(chan int)
	// var d []byte
	// s := fmt.Sprintf("源服务器错误,[状态码]:%v, body:%s", 200, d)
	// fmt.Println(s)
	go tfor()
	// go tfor()
	// go tfor()
	// tfor()
	<-t
}

func tfor() {
	for index := 0; ; index++ {
		// time.Sleep(1 * time.Nanosecond)
	}
}
