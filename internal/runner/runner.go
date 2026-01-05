package runner

import (
	"io"
	"os"
	"os/exec"
)

type Result struct {
	ExitCode int
}

func Run(args []string, logWriter io.Writer) (*Result, error) {
	if len(args) == 0 {
		return nil, exec.ErrNotFound
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin

	if logWriter != nil {
		cmd.Stdout = io.MultiWriter(os.Stdout, logWriter)
		cmd.Stderr = io.MultiWriter(os.Stderr, logWriter)
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, err
		}
	}

	return &Result{ExitCode: exitCode}, nil
}
