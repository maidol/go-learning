package main

import (
	"fmt"
)

var global *int

func f() {
	fmt.Println("global", global)
	var x int
	fmt.Println(&x)
	x = 1
	global = &x
	fmt.Println("global", global)

}

func g() {
	y := new(int)
	*y = 1
}

type myInt int

// AB
func (mi *myInt) AB() {
	fmt.Println(mi)
	fmt.Println(*mi)
}

// ABC :
func ABC() {

}

type T struct {
	Name string
}

func main() {
	a := 1
	ap := &a
	b := 1
	bp := &a
	fmt.Printf("%p, %p\n", &a, &b)
	fmt.Printf("%p, %p\n", ap, bp)
	c := &T{}
	fmt.Printf("%p\n", c)
	c = &T{}
	fmt.Printf("%p\n", c)
	d := T{}
	fmt.Printf("%p\n", &d)
	// a = "afffffffffffffffffffffffffffffffffffff"
	// fmt.Printf("%p\n", &a)
	// f()
	// p := new(int)
	// fmt.Println(p)
	// *p = 123
	// fmt.Println(p)
	// fmt.Println(*p)

	// m := make(map[string]int)
	// v, ok := m["a"]
	// fmt.Println(v, ok)

	// // medals := make([]string, 10)
	// medals := []string{"gole", "silver", "bronze"}
	// medals[0] = "a"
	// medals[1] = "a"
	// medals[2] = "a"
	// fmt.Println(medals)
	// medals = nil

	// var m myInt = 1
	// m.AB()
	// c := tempconv.Celsius(2.1)
	// fmt.Println(c.String())

	s := `你好`
	fmt.Println(len(s))
	fmt.Println(s[2])
	s1 := s[3:6]
	fmt.Println(s1)
	fmt.Println([]byte(s))
	fmt.Printf("%b\n", []byte(s))

	tt(&TT{})

	mt := TT{}
	fmt.Println(mt == TT{})
}

func tt(ta TTA) {
	tt, ok := ta.(*TT)
	fmt.Println(ok)
	ta.s()
	tt.s()
}

type TTA interface {
	s()
}

type TT struct{}

// func (t *T) s() {
// 	fmt.Println("T")
// }

func (t TT) s() {
	fmt.Println("TT")
}
