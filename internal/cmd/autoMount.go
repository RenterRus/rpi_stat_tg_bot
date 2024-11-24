package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func (c *CMD) Auto() string {
	md, err := c.finder.FindMD()
	if err != nil {
		return fmt.Sprintf("auto(md) error: %s", err.Error())
	}

	mountCMD, mountARG := c.informer.CMDMount(md)
	chmodCMD, chmodARG := c.informer.CMDChmod()

	ARGstr := strings.Builder{}
	ARGstr.WriteString(mountCMD)
	ARGstr.WriteString(" ")
	for _, v := range mountARG {
		ARGstr.WriteString(v)
		ARGstr.WriteString(" ")
	}

	cmd := exec.Command("/bin/sh", "-c", ARGstr.String())
	if err = cmd.Run(); err != nil {
		return fmt.Sprintf("auto(mount) error: %s\ncommand: %s", err.Error(), ARGstr.String())
	}

	ARGstr.Reset()
	ARGstr.WriteString(chmodCMD)
	ARGstr.WriteString(" ")
	for _, v := range chmodARG {
		ARGstr.WriteString(v)
		ARGstr.WriteString(" ")
	}

	cmd = exec.Command("/bin/sh", "-c", ARGstr.String())
	if err = cmd.Run(); err != nil {
		return fmt.Sprintf("auto(chmod) error: %s\ncommand: %s", err.Error(), ARGstr.String())
	}

	return "auto-connection attempt completed"
}
