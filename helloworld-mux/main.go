package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/gorilla/mux"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func simpleMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", MyHandler)

	r.HandleFunc("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/v1"))
	})
	r.HandleFunc("/v1/t", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/v1/t"))
	})
	r.PathPrefix("/v3/pro/bookmark").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("prefix /v3/pro/bookmark"))
	})
	r.PathPrefix("/v3/pro/book").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("prefix /v3/pro/book"))
	})
	r.PathPrefix("/v3/pro").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("prefix /v3/pro"))
	})

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}
