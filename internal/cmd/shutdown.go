package cmd

import (
	"fmt"
	"os/exec"
)

func (*CMD) Shutdown() (string, bool) {
	if err := exec.Command("/bin/sh", "-c", "sudo shutdown -h 3").Run(); err != nil {
		return fmt.Sprintf("restart error: %s", err.Error()), false
	}

	return "Shutdown is planing after 5 minutes", true
}
