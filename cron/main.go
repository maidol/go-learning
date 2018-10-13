package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	g := sync.WaitGroup{}
	g.Add(1)
	go run(&g)
	g.Wait()
}

func run(g *sync.WaitGroup) {
	ch := make(chan interface{})
	t1 := time.NewTimer(5 * time.Second)
	t2 := time.NewTimer(10 * time.Second)
	for {
		select {
		case m := <-ch:
			fmt.Println(m)
			g.Done()
		case <-t1.C:
			fmt.Println("5 seconds")
			cmd("date", "-R")
			t1.Reset(5 * time.Second)
		case <-t2.C:
			fmt.Println("10 seconds")
			cmd("date", "-R")
			t2.Reset(10 * time.Second)
		}
	}
}

func cmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout //
	cmd.Run()
}
