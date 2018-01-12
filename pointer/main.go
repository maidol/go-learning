package main

import (
	"fmt"
)

func main() {
	// dd := map[string]string{"d1": "dv"}
	// m := make(map[string]string)
	// pm := &m
	// m["k1"] = "v11"
	// (*pm)["k1"] = "v111"
	// *pm = map[string]string{"t1": "tv"}
	// m["t1"] = "tv1"
	// m = make(map[string]string)
	// pm = &dd
	// fmt.Println(&m, pm, *pm)
	// fmt.Printf("%p %p %p %p\n", m, *pm, &dd, dd)

	t1 := T{name: "king"}
	p1 := &t1
	p2 := p1
	fmt.Println(t1, p1, p2)
	t1 = T{name: "tom"}
	fmt.Println(t1, p1, p2)
	*p1 = T{name: "lili"}
	fmt.Println(t1, p1, p2)
	(*p2).name = "lili1"
	fmt.Println(t1, p1, p2)
	p1 = &T{name: "tim"}
	fmt.Println(t1, p1, p2)
}

type T struct {
	name string
}
