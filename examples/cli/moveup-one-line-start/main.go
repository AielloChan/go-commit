package main

import (
	"os"
	"time"
)

func main() {
	os.Stdout.WriteString("000")
	os.Stdout.WriteString("\n111")
	os.Stdout.WriteString("\n222222")
	time.Sleep(1e3 * time.Millisecond)
	// move cursor to the beginning of the previous line
	os.Stdout.WriteString("\033[F")
	os.Stdout.WriteString("333\n")
}
