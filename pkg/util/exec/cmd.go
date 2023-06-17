package exec

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Cmd(dir, name string, args ...string) error {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(name, args[:]...) // #nosec
	if len(dir) >0 {
		cmd.Dir = dir
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func Output(name string, args ...string) ([]byte, error) {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(name, args[:]...) // #nosec
	return cmd.CombinedOutput()
}

func RunSimpleCmd(cmd string) (string, error) {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	result, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput() // #nosec
	return string(result), err
}

func RunBashCmd(cmd string) (string, error) {
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	result, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput() // #nosec
	return string(result), err
}

func BashEval(cmd string) string {
	out, _ := RunBashCmd(cmd)
	return strings.Trim(out, "\n\t")
}

func Eval(cmd string) string {
	out, _ := RunSimpleCmd(cmd)
	return strings.Trim(out, "\n\t")
}

func CheckCmdIsExist(cmd string) (string, bool) {
	cmd = fmt.Sprintf("type %s", cmd)
	out, err := RunSimpleCmd(cmd)
	if err != nil {
		return "", false
	}

	outSlice := strings.Split(out, "is")
	last := outSlice[len(outSlice)-1]

	if last != "" && !strings.Contains(last, "not found") {
		return strings.TrimSpace(last), true
	}
	return "", false
}
