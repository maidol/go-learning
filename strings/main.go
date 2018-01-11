package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Replace("/gateway/user/add", "/gateway/", "/", 1))
	fmt.Println("/v5/epaperwork/getReceiveBookchapters/", strings.Trim("/v5/epaperwork/getReceiveBookchapters/", "/"), strings.TrimRight("/v5/epaperwork/getReceiveBookchapters/", "/"))
	fmt.Println(len(strings.Split(strings.Trim("/v5/epaperwork/getReceiveBookchapters/", "/"), "/")), len(strings.Split(strings.Trim("/v5/epaperwork/", "/"), "/")))
}
