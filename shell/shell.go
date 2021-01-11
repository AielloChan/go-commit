package shell

import (
	"os/exec"
	"strings"
)

func ProcessShell(cmd string) string {
	cmdResult := exec.Command("/bin/sh", "-c", cmd)
	if res, err := cmdResult.Output(); err != nil {
		panic("Can't run shell '" + cmd + "' with error: " + err.Error())
	} else {
		output := string(res)
		if strings.HasSuffix(output, "\n") {
			return output[:len(output)-1]
		}
		return output
	}
}
