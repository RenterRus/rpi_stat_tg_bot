package finder

import (
	"fmt"
	"os"
	"strings"
)

func (k *RealFinder) getMD() (string, error) {
	devices, err := os.ReadDir("/dev")
	if err != nil {
		return "", fmt.Errorf("getMD (ReadDir): %w", err)
	}

	for _, device := range devices {
		// Игнорируем неинтересующие нас файлы
		if device.IsDir() || device.Name() == "stdin" || device.Name() == "stdout" || device.Name() == "stderr" {
			continue
		}

		if strings.Contains(device.Name(), k.fileSearch) {
			return device.Name(), nil
		}
	}

	return "", ErrNoFound
}

func (k *RealFinder) FindMD() (string, error) {
	device, err := k.getMD()
	if err != nil {
		return "", fmt.Errorf("FindMD(): %w", err)
	}

	return device, nil
}
