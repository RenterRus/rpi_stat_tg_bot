package cmd

import (
	"fmt"
	"syscall"
)

func (c *CMD) Auto() string {
	md, err := c.finder.FindMD()
	if err != nil {
		return fmt.Sprintf("auto(md) error: %s", err.Error())
	}

	mountCMD, mountARG := c.informer.CMDMount(md)
	chmodCMD, chmodARG := c.informer.CMDChmod()

	err = syscall.Exec(mountCMD, mountARG, nil)
	if err != nil {
		return fmt.Sprintf("auto(mount) error: %s", err.Error())
	}

	err = syscall.Exec(chmodCMD, chmodARG, nil)
	if err != nil {
		return fmt.Sprintf("auto(chmod) error: %s", err.Error())
	}

	return "auto-connection attempt completed"
}
