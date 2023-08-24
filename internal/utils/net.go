package utils

import (
	"net"
)

// GlobalIP4 get host global unicast ipv4 address.
func GlobalIP4() (string, error) {
	var (
		err       error
		ip        net.IP
		addresses []net.Addr
	)

	if addresses, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)

		if !ok || !ipNet.IP.IsGlobalUnicast() {
			continue
		}

		if ipv4 := ipNet.IP.To4(); ipv4 != nil {
			ip = ipv4
			break
		}
	}

	if ip != nil {
		return ip.String(), nil
	}

	return "", nil
}
