package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Fprintln(os.Stderr, "Empty command list")
		return 1
	}

	if err := env.Apply(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set environmetns: %s\n", err)
		return 1
	}

	root := cmd[0]
	command := exec.Command(root)
	if len(cmd) > 1 {
		command.Args = append(command.Args, cmd[1:]...)
	}

	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		fmt.Fprintf(os.Stderr, "Faild to run command: %s\n", err)
		return 1
	}

	return 0
}
