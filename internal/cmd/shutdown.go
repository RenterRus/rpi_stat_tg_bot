package cmd

import (
	"fmt"
	"syscall"
)

func (*CMD) Shutdown() (string, bool) {
	err := syscall.Exec("sudo shutdown", []string{"-h 5"}, nil)
	if err != nil {
		return fmt.Sprintf("shutdown error: %s", err.Error()), false
	}

	return "Shutdown is planing after 5 minutes", true
}
