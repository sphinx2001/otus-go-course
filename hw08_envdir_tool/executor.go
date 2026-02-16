package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// #nosec G204
	command := exec.Command(cmd[0], cmd[1:]...)
	preparedEnvs := PrepareEnv(env)

	command.Env = preparedEnvs
	var out strings.Builder
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		fmt.Println("Exit code:", err)
	}
	fmt.Println(out.String())
	return command.ProcessState.ExitCode()
}
