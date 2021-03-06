package main

import (
	"net/http"
	// _ "net/http/pprof"
	"runtime"

	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

func main() {
	// go http.ListenAndServe(":8888", nil)
	runtime.GOMAXPROCS(4)
	fwd, _ := forward.New()
	redirect := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.URL = testutils.ParseURI("http://192.168.2.101:8708")

		fwd.ServeHTTP(w, r)
	})
	s := &http.Server{
		Addr:    ":8080",
		Handler: redirect,
	}
	s.ListenAndServe()
}
