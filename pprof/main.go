package main

import "log"
import "net/http"
import _ "net/http/pprof"

func main() {
	ch := make(chan int)
	go func() {
		log.Fatalln(http.ListenAndServe("localhost:6060", nil))
	}()
	<-ch
}
