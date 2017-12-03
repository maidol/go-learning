package main

import "fmt"

type my []string

type myi interface{}

func main() {
	defer func() {
		fmt.Println("test recover")
		e, ok := recover().(my)
		if ok {
			fmt.Println(e[0])
		}

	}()
	a := my{"a", "b"}
	var myinter myi = a
	b := myinter
	fmt.Println(a, b, myinter)
	p()
}

func p() {
	panic(my{"a", "b"})
}
