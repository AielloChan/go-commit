package tools

import (
	"bytes"
	"html/template"
	"strings"

	"yuanling.com/go-commit/model"
	"yuanling.com/go-commit/shell"
)

func ProcessTpl(tplStr string, store *model.Store) string {
	var preview bytes.Buffer
	tpl, err := template.New("preview").Parse(tplStr)
	if err != nil {
		panic("format preview string error")
	}
	err = tpl.Execute(&preview, &store)
	if err != nil {
		panic("Parse template '" + tplStr + "' error: " + err.Error())
	}
	return strings.ReplaceAll(preview.String(), "\\n", "\n")
}

func GetString(s string, store *model.Store) string {
	switch {
	case strings.HasPrefix(s, "#!"):
		cmd := s[len("#!"):]
		return shell.ProcessShell(cmd)
	case strings.HasPrefix(s, "#$"):
		tpl := s[len("#$"):]
		return ProcessTpl(tpl, store)
	case strings.HasPrefix(s, "#!$"):
		mixedStr := s[len("#!$"):]
		shellRes := shell.ProcessShell(mixedStr)
		return ProcessTpl(shellRes, store)
	case strings.HasPrefix(s, "#$!"):
		mixedStr := s[len("#$!"):]
		tplRes := ProcessTpl(mixedStr, store)
		return shell.ProcessShell(tplRes)
	default:
		return s
	}
}
