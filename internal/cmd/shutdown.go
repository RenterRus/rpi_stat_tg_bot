package cmd

import (
	"fmt"
	"os/exec"
)

func (c *CMD) RestartBot(name string) (string, bool) {
	if err := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo systemctl restart %s.service", name)).Run(); err != nil {
		return fmt.Sprintf("restart error: %s", err.Error()), false
	}

	return "Бот перезапускается", true
}
