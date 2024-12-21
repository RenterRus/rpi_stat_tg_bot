package cmd

import (
	"fmt"
	"os/exec"
)

func (c *CMD) Restart() (string, bool) {
	if err := exec.Command("/bin/sh", "-c", "sudo reboot now").Run(); err != nil {
		return fmt.Sprintf("reboot error: %s", err.Error()), false
	}

	return "Перезапуск. Не забудьте подключить RAID при включении", true
}
