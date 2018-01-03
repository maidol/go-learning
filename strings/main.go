package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Replace("/gateway/user/add", "/gateway/", "/", 1))
}
