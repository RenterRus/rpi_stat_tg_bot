package informer

import "fmt"

func (k *KekInformer) CMDMount(md string) (string, []string) {
	return "sudo mount", []string{fmt.Sprintf("/dev/%s", md), fmt.Sprintf("/home/%s", k.root_user)}
}
