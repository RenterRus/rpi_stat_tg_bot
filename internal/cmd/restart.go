package cmd

import (
	"fmt"
	"syscall"
)

func (*CMD) Restart() (string, bool) {
	err := syscall.Exec("sudo shutdown", []string{"-r 1"}, nil)
	if err != nil {
		return fmt.Sprintf("restart error: %s", err.Error()), false
	}

	return "Restart is planing after minutes", true
}
