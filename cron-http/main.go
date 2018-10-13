package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	port := os.Args[1]
	http.HandleFunc("/", crontask)
	fmt.Println("GoCron listening on " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

var taskMap = make(map[string]*Task)
var m = sync.RWMutex{}

type Task struct {
	Id       string   `json:"id"`
	Interval int      `json:"interval"`
	Cmd      string   `json:"cmd"`
	Args     []string `json:"args"`
	stop     chan interface{}
}

type cronResponse struct {
	Ok    bool   `json:"ok"`
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func crontask(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	fmt.Println(method)
	switch method {
	case "POST":
		createTask(w, r)
	case "DELETE":
		removeTask(w, r)
	}

}

func createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	d, _ := ioutil.ReadAll(r.Body)
	v := &Task{}
	json.Unmarshal(d, v)
	fmt.Println(v)
	m.Lock()
	defer m.Unlock()
	_, ok := taskMap[v.Id]
	if ok {
		cr := cronResponse{true, v.Id, "The task " + v.Id + " already exists."}
		crd, _ := json.Marshal(&cr)
		w.WriteHeader(409)
		w.Write(crd)
		return
	}
	taskMap[v.Id] = v
	go run(v)

	cr := cronResponse{true, v.Id, ""}
	crd, _ := json.Marshal(&cr)
	w.Write(crd)
}

func run(t *Task) {
	t.stop = make(chan interface{})
	t1 := time.NewTimer(time.Duration(t.Interval/1000) * time.Second)
	for {
		select {
		case <-t.stop:
			fmt.Println("stop task")
			return
		case <-t1.C:
			// fmt.Println("5 seconds")
			cmd(t.Cmd, t.Args...)
			t1.Reset(5 * time.Second)
		}
	}
}

func cmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout //
	cmd.Run()
}

func removeTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	ok := remove(id)
	var rd cronResponse
	var crd []byte
	if !ok {
		rd = cronResponse{false, id, "The task " + id + " is not found."}
		crd, _ = json.Marshal(&rd)
		w.WriteHeader(404)
		w.Write(crd)
		return
	}
	rd = cronResponse{true, id, ""}
	crd, _ = json.Marshal(&rd)
	w.Write(crd)
}

func remove(id string) bool {
	m.Lock()
	defer m.Unlock()
	t, ok := taskMap[id]
	if !ok {
		return false
	}
	t.stop <- 1
	delete(taskMap, id)
	return true
}
