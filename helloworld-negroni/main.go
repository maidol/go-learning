package main

import (
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func MW1(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

func main() {
	runtime.GOMAXPROCS(4)
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi, negroni hello world!"))
	})

	n := negroni.Classic()
	// n.Use(negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 	next.ServeHTTP(rw, r)
	// }))
	n.UseFunc(MW1)
	// router goes last
	n.UseHandler(router)

	n.Run(":8706")
}
