package finder

import "net"

type Finder interface {
	FindMD() (string, error)
	FindIP() ([]net.IP, error)
}
