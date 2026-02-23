package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// #nosec G204
	command := exec.Command(cmd[0], cmd[1:]...)
	preparedEnvs := PrepareEnv(env)

	command.Env = preparedEnvs
	command.Stdin = os.Stdin   // Пробрасываем ввод
	command.Stdout = os.Stdout // Пробрасываем вывод
	command.Stderr = os.Stderr // Пробрасываем ошибки
	err := command.Run()
	if err != nil {
		fmt.Println("Exit code:", err)
	}
	return command.ProcessState.ExitCode()
}
