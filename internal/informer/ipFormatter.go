package informer

import (
	"fmt"
	"strings"
)

func (k *RealInformer) IPFormatter() string {
	result := strings.Builder{}

	result.WriteString("Host ip: ")
	ips, err := k.finder.FindIP()
	if err != nil {
		return fmt.Errorf("error into search ip. Reason: %w", err).Error()
	}

	for i, v := range ips {
		result.WriteString(v.String())
		if i > 0 && i < len(ips)-1 {
			result.WriteString(", ")
		}
	}
	result.WriteString(" ")
	result.WriteString("\n")
	return result.String()
}
