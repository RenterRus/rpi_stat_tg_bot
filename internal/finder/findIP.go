package finder

import (
	"fmt"
	"net"
)

func (*KekFinder) getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}

func (k *KekFinder) FindIP() ([]net.IP, error) {
	ips, err := k.getLocalIPs()
	if err != nil {
		return nil, fmt.Errorf("FindIP(): %w", err)
	}

	return ips, nil
}
