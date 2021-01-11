package main

import (
	"os"
	"time"
)

func main() {
	os.Stdout.WriteString("000")
	os.Stdout.WriteString("\n111")
	os.Stdout.WriteString("\n222")
	time.Sleep(1e3 * time.Millisecond)
	os.Stdout.WriteString("\r\x1b[K")
	os.Stdout.WriteString("333\n")
}
