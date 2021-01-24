package tools

import (
	"os"

	"github.com/fatih/color"
)

func PrintError(msg string) {
	color.New(color.FgRed, color.Bold).Fprint(os.Stdout, msg)
}
