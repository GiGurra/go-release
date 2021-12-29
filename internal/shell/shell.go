package shell

import (
	"os/exec"
	"strings"
)

func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func RunCommand(program string, args ...string) string {

	if !CommandExists(program) {
		panic("k8s not on path!")
	}

	fullCommand := program + strings.Join(args, " ")

	cmd := exec.Command(program, args...)

	outputBytes, err := cmd.Output()

	if err != nil {
		exitErr, isExitErr := err.(*exec.ExitError)
		if isExitErr {
			panic(`command "` + fullCommand + `" failed with error: ` + exitErr.Error() + ", stderr: " + string(exitErr.Stderr))
		} else {
			panic(`command "` + fullCommand + `" failed with error: ` + err.Error())
		}
	}

	return strings.TrimSpace(string(outputBytes))
}
