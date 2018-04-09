// package main

// /*
// #include <stdio.h>
// #include <stdlib.h>
// */
// import "C"
// import "unsafe"

// func main() {
// 	cstr := C.CString("HELLO, WORLD")
// 	C.puts(cstr)
// 	C.free(unsafe.Pointer(cstr))
// }

// sample: https://github.com/hyper0x/go_command_tutorial/blob/master/0.13.md
package main

import (
	"fmt"
	cgolib "go-learning/cgo/lib"
)

func main() {
	input := float32(2.33)
	output, err := cgolib.Sqrt(input)
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
	fmt.Printf("The square root of %f is %f.\n", input, output)

	cgolib.Print("ABC\n")

	cgolib.CallCFunc()
}
