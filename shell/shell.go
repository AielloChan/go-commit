package shell

import (
	"errors"
	"os/exec"
	"strings"
)

func ProcessShell(cmd string) (string, error) {
	cmdResult := exec.Command("/bin/sh", "-c", cmd)
	res, err := cmdResult.Output()
	if err != nil {
		return "", errors.New("Excute command '" + cmd + "' error: " + err.Error())
	}
	output := string(res)
	if strings.HasSuffix(output, "\n") {
		return output[:len(output)-1], nil
	}
	return output, nil
}
