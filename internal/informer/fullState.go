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

	return fmt.Sprintf("RAID size:\nTotal: %d gb(round int)\nFree: %d gb(round int)\nAvailable: %d gb(round int)\nUsed: %d gb(round int)\n",
			int(((stat.Blocks*uint64(stat.Bsize)/1024)/1024)/1024),
			int(((stat.Bfree*uint64(stat.Bsize)/1024)/1024)/1024),
			int(((stat.Bavail*uint64(stat.Bsize)/1024)/1024)/1024),
			int((((stat.Blocks*uint64(stat.Bsize))-(stat.Bfree*uint64(stat.Bsize)))/1024)/1024)/1024),
		nil
}
