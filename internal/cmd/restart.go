package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
)

func (c *CMD) Restart() (string, bool) {
	if err := exec.Command("/bin/sh", "-c", "sudo shutdown -r +"+strconv.Itoa(c.ttp)).Run(); err != nil {
		return fmt.Sprintf("reboot error: %s", err.Error()), false
	}

	return "Restart is planing after " + strconv.Itoa(c.ttp) + " minutes. But the bot shutdown right now", true
}
