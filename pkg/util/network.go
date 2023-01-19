package util

import (
	"fmt"
	"net"
)

func GetAvailableEphemeralPort(iface string) (int, error) {
	address := fmt.Sprintf("%s:0", iface)

	listen, err := net.Listen("tcp", address)
	if err != nil {
		return 0, err
	}
	defer listen.Close()

	tcp, _ := listen.Addr().(*net.TCPAddr)

	return tcp.Port, nil
}
