package informer

import (
	"fmt"
	"syscall"
)

func (k *KekInformer) FullState() (string, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/home/"+k.root_user, &stat); err != nil {
		return "", fmt.Errorf("syscall.Statfs: %w", err)
	}

	return fmt.Sprintf("Size\nTotal: %d\nFree: %d\nAvailable: %d\nUsed: %d\n",
			stat.Blocks*uint64(stat.Bsize),
			stat.Bfree*uint64(stat.Bsize),
			stat.Bavail*uint64(stat.Bsize),
			(stat.Blocks*uint64(stat.Bsize))-(stat.Bfree*uint64(stat.Bsize))),
		nil
}
