package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("command receives correct environment", func(t *testing.T) {
		env, err := ReadDir("./testdata/env/")
		require.NoError(t, err)

		cmd := []string{"env"}

		var out bytes.Buffer
		comand := exec.Command(cmd[0], cmd[1:]...)

		cmdEnv := []string{}
		for k, v := range env {
			if !v.NeedRemove {
				cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", k, v.Value))
			}
		}
		comand.Env = cmdEnv
		comand.Stdout = &out
		comand.Stderr = &out

		err = comand.Run()
		require.NoError(t, err)

		output := out.String()

		require.NotContains(t, output, "UNSET=")
		require.Contains(t, output, "BAR=bar")
		require.Contains(t, output, "FOO=   foo\nwith new line")
		require.Contains(t, output, `HELLO="hello"`)
		require.Contains(t, output, "EMPTY=")
	})
}
