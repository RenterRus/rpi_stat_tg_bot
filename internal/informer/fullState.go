package informer

import (
	"fmt"
	"syscall"
)

func (k *RealInformer) FullState() (string, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/home/"+k.root_user, &stat); err != nil {
		return "", fmt.Errorf("syscall.Statfs: %w", err)
	}

	return fmt.Sprintf("RAID size:\nTotal: %.2f gb\nFree: %.2f gb\nAvailable: %.2f gb\nUsed: %.2f gb\n",
		(float64(stat.Blocks*uint64(stat.Bsize)/1024)/1024)/1024,
		(float64(stat.Bfree*uint64(stat.Bsize)/1024)/1024)/1024,
		(float64(stat.Bavail*uint64(stat.Bsize)/1024)/1024)/1024,
		(float64((stat.Blocks*uint64(stat.Bsize)-stat.Bfree*uint64(stat.Bsize))/1024)/1024)/1024), nil
}
