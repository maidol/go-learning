package main

import "fmt"
import "helloworld/comm"

type person struct {
	Name string
	Age  int
}

func main() {
	fmt.Printf("hello world\n")
	comm.Add()
	var trace int32 = 1
	var trace1 int32 = 1
	fmt.Println(trace + trace1)

	p := person{}
	p.Name = "mark"
	p.Age = 18
	p.Print()
	(&p).Print()
	fmt.Println(p)

	c := make(chan bool)
	go func() {
		fmt.Println("go routine")
		c <- true
	}()
	<-c
}

func (p person) Print() {
	fmt.Println("print")
}
