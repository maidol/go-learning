package main

import (
	"bytes"
	"fmt"
)

func main() {
	arr := [5]string{"a", "b", "c", "d", "e"}
	aslice := arr[2:4]
	bslice := arr[:4]
	cslice := arr[4:]
	dslice := arr[:2]
	fmt.Printf("arr {\"a\", \"b\", \"c\", \"d\", \"e\"} len=%d,cap=%d\n", len(arr), cap(arr))
	fmt.Printf("aslice arr[2:4] len=%d,cap=%d\n", len(aslice), cap(aslice))
	fmt.Printf("bslice arr[:4] len=%d,cap=%d\n", len(bslice), cap(bslice))
	fmt.Printf("cslice arr[4:] len=%d,cap=%d\n", len(cslice), cap(cslice))
	fmt.Printf("dslice arr[:2] len=%d,cap=%d\n", len(dslice), cap(dslice))

	buf := bytes.Buffer{}
	buf.WriteString("hellohelloaaahellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello")
	fmt.Printf("buf: %s\n", buf.Bytes())
	b := make([]byte, 10)
	buf.Read(b)
	c := make([]byte, 8)
	buf.Read(c)
	fmt.Printf("read b:%s,c:%s\n", b, c)
}
