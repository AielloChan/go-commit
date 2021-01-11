package main

import (
	"fmt"
	"regexp"
)

func main() {
	m1 := regexp.MustCompile(`(?![^\\])(\n)`)

	fmt.Println(m1.ReplaceAllString("123\n456", "#"))
	fmt.Println(m1.ReplaceAllString("123\\n456", "#"))
}
