package main

import (
	"fmt"
)

func main() {
	defer func() {
		r := recover()
		fmt.Println(r)
	}()
	testPanic()
}

func testPanic() {
	panic("test panic")
}
