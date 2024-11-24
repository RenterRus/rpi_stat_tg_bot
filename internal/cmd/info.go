package cmd

import (
	"strings"
)

func (c *CMD) Info() (string, string) {
	message := strings.Builder{}
	m := ""
	var err error
	cmd := ""
	if m, cmd, err = c.informer.Basic(); err != nil {
		message.WriteString("Error branch\n")
		message.WriteString("Reason: ")
		message.WriteString(err.Error())
		message.WriteString("\n")

		_, err := c.finder.FindMD()
		if err != nil {
			message.WriteString("Finding md error: ")
			message.WriteString(err.Error())
			message.WriteString("\n")
		}

		message.WriteString("\n")

		ips, err := c.finder.FindIP()
		message.WriteString("Hosting: ")
		if err != nil {
			message.WriteString("Finding ip error: ")
			message.WriteString(err.Error())
			message.WriteString("\n")
		} else {
			for _, v := range ips {
				message.WriteString(v.String())
				message.WriteString("\n")
			}
		}

		return message.String(), cmd
	}

	return m, cmd
}
