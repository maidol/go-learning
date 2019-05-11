package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
)

func FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
