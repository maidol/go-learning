package main

import (
	"log"
	"net/http"

	// "github.com/gorilla/mux"
	"github.com/roobre/gorilla-mux"
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
	// r.
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}
