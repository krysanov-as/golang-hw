package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	baseMainError   = 1
	commandNotFound = 127
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 || cmd[0] == "" {
		return commandNotFound
	}
	// #nosec
	comand := exec.Command(cmd[0], cmd[1:]...)
	for k, v := range env {
		os.Unsetenv(k)
		if !v.NeedRemove {
			os.Setenv(k, v.Value)
		}
	}
	comand.Env = os.Environ()
	comand.Stdin = os.Stdin
	comand.Stdout = os.Stdout
	comand.Stderr = os.Stderr
	err := comand.Run()
	if err != nil {
		var errExitCode *exec.ExitError
		if errors.As(err, &errExitCode) {
			return errExitCode.ProcessState.ExitCode()
		}
		return baseMainError
	}

	return comand.ProcessState.ExitCode()
}
