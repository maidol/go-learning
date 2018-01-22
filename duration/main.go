package main

import (
	"fmt"
	"time"
)

func main() {
	d := 1 * time.Second
	fmt.Printf("1 * time.Second = int64(%v)\n", int64(d))
}
