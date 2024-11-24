package cmd

import (
	"fmt"
	"os/exec"
)

func (c *CMD) Update() string {
	if err := exec.Command("/bin/sh", "-c", "sudo rm main").Run(); err != nil {
		return fmt.Sprintf("update error: %s", err.Error())
	}

	if err := exec.Command("/bin/sh", "-c", "git pull").Run(); err != nil {
		return fmt.Sprintf("update error: %s", err.Error())
	}

	if err := exec.Command("/bin/sh", "-c", "go build cmd/main.go").Run(); err != nil {
		return fmt.Sprintf("update error: %s", err.Error())
	}

	if err := exec.Command("/bin/sh", "-c", "sudo systemctl reboot runbot.service").Run(); err != nil {
		return fmt.Sprintf("update error: %s", err.Error())
	}

	return "Attempt to update is running. Maybe ruquried reboot"
}
