package informer

import "fmt"

func (k *KekInformer) CMDChmod() (string, []string) {
	return "sudo chown", []string{k.root_user, fmt.Sprintf("/home/%s/raid/", k.root_user)}
}
