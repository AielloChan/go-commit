package main

import (
	"os"

	"github.com/fatih/color"
)

func main() {
	color.New(color.FgBlue, color.Bold).Fprintf(os.Stdout, "222\n")
}
