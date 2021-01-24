package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmdResult := exec.Command("/bin/sh", "-c", "git status")
	_, err := cmdResult.Output()
	if err != nil {
		fmt.Println("not a git repository")
	} else {
		fmt.Println("is a git repository")
	}
}
