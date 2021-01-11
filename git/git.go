package git

import "yuanling.com/go-commit/shell"

func Commit(message string) {
	shell.ProcessShell("git commit -m " + message)
}
