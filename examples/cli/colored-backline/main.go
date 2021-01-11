package main

import (
	"os"
	"time"

	"github.com/fatih/color"
)

func main() {
	os.Stdout.WriteString("000")
	os.Stdout.WriteString("\n111")
	color.New(color.FgBlue, color.Bold).Fprintf(os.Stdout, "\n22222")
	time.Sleep(1e3 * time.Millisecond)
	os.Stdout.WriteString("\r\x1b[K")
	os.Stdout.WriteString("333\n")
}
