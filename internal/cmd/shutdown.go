package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
)

func (c *CMD) Shutdown() (string, bool) {
	if err := exec.Command("/bin/sh", "-c", "sudo shutdown -h "+strconv.Itoa(c.ttp)).Run(); err != nil {
		return fmt.Sprintf("restart error: %s", err.Error()), false
	}

	return "Shutdown is planing after " + strconv.Itoa(c.ttp) + " minutes. But the bot shutdown right now", true
}
