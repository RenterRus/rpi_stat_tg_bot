package cmd

import (
	"fmt"
	"os/exec"
)

func (c *CMD) Update() (string, bool) {
	if err := exec.Command("/bin/sh", "-c", "sudo rm main && git pull && go build cmd/main.go && sudo systemctl stop runbot.service && sudo systemctl start runbot.service && sudo systemctl enable runbot.service && sudo systemctl status runbot.service").Run(); err != nil {
		return fmt.Sprintf("reboot error: %s", err.Error()), false
	}

	return "Restart is run", true
}
