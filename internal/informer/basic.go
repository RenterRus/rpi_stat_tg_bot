package informer

import (
	"fmt"
	"log"
	"strings"
)

func (k *RealInformer) Basic() (string, string, error) {
	basic := strings.Builder{}

	basic.WriteString(k.IPFormatter())

	md, err := k.finder.FindMD()
	if err != nil {
		return "", "", fmt.Errorf("Basic (MD): %w", err)
	}
	log.Println(md)
	basic.WriteString("\nFinding storage: ")
	basic.WriteString(md)

	cmd := strings.Builder{}

	cmd.WriteString("sudo mount /dev/")
	cmd.WriteString(md)
	cmd.WriteString(" /home/")
	cmd.WriteString(k.root_user)
	cmd.WriteString(" && ")
	cmd.WriteString("sudo chown -R")
	cmd.WriteString(k.root_user)
	cmd.WriteString(" /home/")
	cmd.WriteString(k.root_user)
	cmd.WriteString("/raid/")

	basic.WriteString("\n")

	full, err := k.FullState()
	if err != nil {
		basic.WriteString(fmt.Sprintf("Full state output error: %s\n", err.Error()))
		return basic.String(), cmd.String(), nil
	}

	basic.WriteString(full)

	return basic.String(), cmd.String(), nil
}
