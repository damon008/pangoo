package util

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"pangoo/pkg/constants"
	"net"
	"strings"
)

func GetIpAddr() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		hlog.Error(err)
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// 192.168.1.20:61085
	ip := strings.Split(localAddr.String(), ":")[0]

	return ip
}

// GetFreePort get a free port.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr(constants.TCP, constants.FreePortAddress)
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP(constants.TCP, addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
