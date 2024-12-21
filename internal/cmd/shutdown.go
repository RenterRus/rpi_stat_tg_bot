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

	return "Выключение запланировано через " + strconv.Itoa(c.ttp) + " минут. Но бот выключен уже сейчас", true
}
