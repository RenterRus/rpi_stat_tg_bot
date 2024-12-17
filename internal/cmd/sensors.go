package cmd

import (
	"fmt"
	"os/exec"
)

// Install and configure
// sudo apt install lm-sensors
// sudo sensors-detect

func (c *CMD) Sensors() string {
	var out []byte
	var err error

	if out, err = exec.Command("/bin/sh", "-c", "sensors").Output(); err != nil {
		return fmt.Sprintf("sensors error: %s", err.Error())
	}

	return string(out)
}
