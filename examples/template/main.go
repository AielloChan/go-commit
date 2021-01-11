package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type Store = map[string]string

func main() {
	// tpl, err := template.New("preview").Parse(`Name: {{.Name}}, \nID: {{.ID}}`)
	// if err != nil {
	// 	panic("format preview string error")
	// }
	store := Store{
		"Name": "Aiello",
		"ID":   "123",
	}
	// tpl.Execute(os.Stdout, store)

	var preview bytes.Buffer
	tpl, err := template.New("preview").Parse(`Name: {{.Name}}, \nID: {{.ID}}`)
	if err != nil {
		panic("format preview string error")
	}
	err = tpl.Execute(&preview, &store)
	if err != nil {
		panic("Parse template error:" + err.Error())
	}
	output := preview.String()
	fmt.Println(strings.ReplaceAll(output, "\\n", "\n"))
}
