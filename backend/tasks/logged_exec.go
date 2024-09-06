package tasks

import (
	"log/slog"
	"os/exec"
	"strings"
)

func Exec(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	slog.Info("executing command", "cmd", cmd)

	err := cmd.Run()
	if err != nil {
		slog.Info("execution error. output:\n" + out.String())
		return err
	}

	slog.Info("finished execution. output:\n" + out.String())
	return nil
}
