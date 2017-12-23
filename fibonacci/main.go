package main

import "fmt"
import "time"

func main() {
	start := time.Now()
	fmt.Println(fibonacci(50))
	fmt.Println(time.Now().Sub(start))
}
func fibonacci(i int) int {
	if i < 2 {
		return i
	}
	return fibonacci(i-2) + fibonacci(i-1)
}
