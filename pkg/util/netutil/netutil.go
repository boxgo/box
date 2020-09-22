package netutil

import (
	"net"
)

// GetGlobalUnicastIP get host global unicast ipv4 address.
func GetGlobalUnicastIP() (ipStr string, err error) {
	if ip, err := getGlobalUnicastIP(); err != nil {
		return ipStr, err
	} else if ip4 := ip.To4(); ip4 != nil {
		return ip4.String(), nil
	}

	return
}

func getGlobalUnicastIP() (ip net.IP, err error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			ip = ipNet.IP
			break
		}
	}

	return
}
