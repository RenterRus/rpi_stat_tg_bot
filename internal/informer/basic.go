package informer

import (
	"fmt"
	"log"
	"strings"
)

func (k *KekInformer) Basic() (string, error) {
	basic := strings.Builder{}

	ips, err := k.finder.FindIP()
	if err != nil {
		return "", fmt.Errorf("Basic (IP): %w", err)
	}
	log.Println(ips)

	basic.WriteString("Running ip: ")
	for i, v := range ips {
		basic.WriteString(v.String())
		if i > 0 && i < len(ips)-1 {
			basic.WriteString(", ")
		}
	}

	md, err := k.finder.FindMD()
	if err != nil {
		return "", fmt.Errorf("Basic (MD): %w", err)
	}
	log.Println(md)
	basic.WriteString("Finding storage: ")
	basic.WriteString(md)

	basic.WriteString("\n------------\n")
	basic.WriteString("Enter this command for fast implement storage into ftp: ")
	basic.WriteString("sudo mount /dev/")
	basic.WriteString(md)
	basic.WriteString(" /home/")
	basic.WriteString(k.root_user)
	basic.WriteString(" && ")
	basic.WriteString("sudo chown ")
	basic.WriteString(k.root_user)
	basic.WriteString(" /home/")
	basic.WriteString(k.root_user)
	basic.WriteString("/raid/")

	basic.WriteString("\n")
	basic.WriteString("\n")

	full, err := k.FullState()
	if err != nil {
		basic.WriteString(fmt.Sprintf("Full state output error: %s\n", err.Error()))
		return basic.String(), nil
	}

	basic.WriteString(full)

	return basic.String(), nil
}
