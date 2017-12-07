package main

import (
	"fmt"
	"net/http"
)

// SingleHost
type SingleHost struct {
	// handler
	handler     http.Handler
	allowedHost string
}

// singlehost server http
func (sh *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == sh.allowedHost {
		sh.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world\n"))
}

type AppendMiddleware struct {
	handler http.Handler
}

func (ad *AppendMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("start => this is a append middleware\n"))
	ad.handler.ServeHTTP(w, r)
	w.Write([]byte("end => this is a append middleware"))
}

type Middleware struct {
	nextHandler   *Middleware
	newestHandler *Middleware
	handler       http.Handler
	postHandler   http.Handler
	name          string
}

func (c *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("begin => " + c.name + " middleware\n"))
	if c.handler != nil {
		c.handler.ServeHTTP(w, r)
	}
	if c.nextHandler != nil {
		c.nextHandler.ServeHTTP(w, r)
	}
	if c.postHandler != nil {
		c.postHandler.ServeHTTP(w, r)
	}
	w.Write([]byte("end => " + c.name + " middleware\n"))
}

type App struct {
	mw *Middleware
}

func (app *App) use(s *Middleware) *Middleware {
	mw := app.mw
	fmt.Println("use middleware => " + s.name)
	if mw.newestHandler == nil {
		mw.nextHandler = s
		mw.newestHandler = s
		return mw
	}
	mw.newestHandler.nextHandler = s
	mw.newestHandler = s
	return mw
}

func New() App {
	return App{
		mw: &Middleware{
			name: "app.md",
		},
	}
}

func main() {
	// single := &SingleHost{
	// 	handler:     http.HandlerFunc(myHandler),
	// 	allowedHost: "example.com",
	// }
	// myad := &AppendMiddleware{
	// 	handler: single,
	// }
	app := New()

	ad1 := &Middleware{
		name:    "ad1",
		handler: http.HandlerFunc(myHandler),
		postHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// r.Context
			w.Write([]byte("postHandler\n"))
		}),
	}
	ad2 := &Middleware{
		name: "ad2",
	}
	ad3 := &Middleware{
		name: "ad3",
	}
	app.use(ad1)
	app.use(ad2)
	app.use(ad3)
	http.ListenAndServe(":8080", app.mw)
}
