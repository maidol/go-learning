package main

import (
	"fmt"
	"net/http"
	// "strings"
	"log"
	"runtime"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()  //解析参数，默认是不会解析的
	// fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	fmt.Fprintf(w, "hi, go-http hello world!") //这个写入到w的是输出到客户端的
}

type MyHandler struct {
}

func (m *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi, go-http hello world!"))
}

func main() {
	runtime.GOMAXPROCS(4)
	http.HandleFunc("/", sayhelloName) //设置访问的路由

	go http.ListenAndServe(":8702", http.Handler(&MyHandler{}))

	err := http.ListenAndServe(":8701", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
