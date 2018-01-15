package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	v := map[string]string{"k1": "v1", "k2": "v2"}
	data, _ := json.Marshal(v)
	fmt.Println(string(data))
}
