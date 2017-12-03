package main

import "fmt"

type mychint chan int

func sum(a []int, c mychint) {
	total := 0
	for _, v := range a {
		total += v
	}
	c <- total // send total to c
}

// func fibonacci(n int, c chan int) {
// 	x, y := 1, 1
// 	for i := 0; i < n; i++ {
// 		c <- x
// 		x, y = y, x+y
// 	}
// 	close(c)
// }

func fibonacci(c, quit chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	// sample1:
	a := []int{7, 2, 8, -9, 4, 0}

	// c := make(chan int)
	c := make(mychint)
	d := make(mychint)
	go sum(a[:len(a)/2], c)
	go sum(a[len(a)/2:], d)
	x, y := <-c, <-d // receive from c, d

	fmt.Println(x, y, x+y)

	// // sample2:
	// c := make(chan int, 1) //修改2为1就报错，修改2为3可以正常运行
	// c <- 1
	// c <- 2
	// fmt.Println(<-c)
	// fmt.Println(<-c)

	// // sample3:
	// c := make(chan int, 10)
	// go fibonacci(cap(c), c)
	// fmt.Println("fibonacci")
	// for i := range c {
	// 	fmt.Println(i)
	// }

	// // sample4:
	// c := make(chan int)
	// quit := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		fmt.Println(<-c)
	// 	}
	// 	quit <- 0
	// }()
	// fibonacci(c, quit)
}
