package informer

import "fmt"

func (k *RealInformer) CMDChmod() (string, []string) {
	return "sudo chown", []string{k.root_user, fmt.Sprintf("/home/%s/raid/", k.root_user)}
}
