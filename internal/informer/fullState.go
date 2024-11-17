package informer

import (
	"fmt"
	"syscall"
)

func (k *KekInformer) FullState() (string, error) {
	var stat syscall.Statfs_t
	md, err := k.finder.FindMD()
	if err != nil {
		return "", fmt.Errorf("FullState: %w", err)
	}

	if err := syscall.Statfs("/dev/"+md, &stat); err != nil {
		return "", fmt.Errorf("syscall.Statfs: %w", err)
	}

	return fmt.Sprintf("Size\nTotal: %d\nFree: %d\nAvailable: %d\nUser: %d\n",
			stat.Blocks*uint64(stat.Bsize),
			stat.Bfree*uint64(stat.Bsize),
			stat.Bavail*uint64(stat.Bsize),
			(stat.Blocks*uint64(stat.Bsize))-(stat.Bfree*uint64(stat.Bsize))),
		nil
}
