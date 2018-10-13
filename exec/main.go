package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("date", "-R")
	cmd.Stdout = os.Stdout //
	cmd.Run()
	// fmt.Println(cmd.Start()) //exec: already started
}
