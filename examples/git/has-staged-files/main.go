package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmdResult := exec.Command("/bin/sh", "-c", "git diff --cached --name-only")
	res, err := cmdResult.Output()
	if err != nil {
		fmt.Println("not a git repository")
	} else {
		pureStr := strings.Trim(string(res), "\n")
		list := strings.Split(pureStr, "\n")
		fmt.Println("staged files: " + strings.Join(list, ","))
	}
}
