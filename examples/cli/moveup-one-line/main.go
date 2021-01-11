package main

import (
	"os"
	"time"
)

func main() {
	os.Stdout.WriteString("000")
	os.Stdout.WriteString("111\n")
	os.Stdout.WriteString("222222\n")
	time.Sleep(1e3 * time.Millisecond)
	// move cursor up one line
	os.Stdout.WriteString("\033[A")
	os.Stdout.WriteString("333\n")
}
