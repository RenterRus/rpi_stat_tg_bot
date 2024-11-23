package cmd

import (
	"fmt"
	"strings"
	"syscall"
)

func (c *CMD) Auto() string {
	md, err := c.finder.FindMD()
	if err != nil {
		return fmt.Sprintf("auto(md) error: %s", err.Error())
	}

	mountCMD, mountARG := c.informer.CMDMount(md)
	chmodCMD, chmodARG := c.informer.CMDChmod()

	ARGstr := strings.Builder{}

	for _, v := range mountARG {
		ARGstr.WriteString(v)
		ARGstr.WriteString(" ")
	}

	err = syscall.Exec(mountCMD, mountARG, nil)
	if err != nil {
		return fmt.Sprintf("auto(mount) error: %s\ncommand: %s %s", err.Error(), mountCMD, ARGstr.String())
	}

	ARGstr.Reset()
	for _, v := range chmodARG {
		ARGstr.WriteString(v)
		ARGstr.WriteString(" ")
	}

	err = syscall.Exec(chmodCMD, chmodARG, nil)
	if err != nil {
		return fmt.Sprintf("auto(chmod) error: %s\ncommand: %s %s", err.Error(), chmodCMD, ARGstr.String())
	}

	return "auto-connection attempt completed"
}
