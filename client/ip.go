package client

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

func parseIp(addr string) (ip net.IP, port int, err error) {
	ipAndPort := strings.Split(addr, ":")
	if len(ipAndPort) != 2 {
		err = errors.New("unkown addr:" + addr)
		return
	}
	ips := strings.Split(ipAndPort[0], ".")
	if len(ips) != 4 {
		err = errors.New("unkown addr:" + addr)
		return
	}
	aa, err := strconv.ParseUint(ips[0], 10, 8)
	if err != nil {
		return
	}
	bb, err := strconv.ParseUint(ips[1], 10, 8)
	if err != nil {
		return
	}
	cc, err := strconv.ParseUint(ips[2], 10, 8)
	if err != nil {
		return
	}
	dd, err := strconv.ParseUint(ips[3], 10, 8)
	if err != nil {
		return
	}
	ee, err := strconv.ParseUint(ipAndPort[1], 10, 16)
	if err != nil {
		return
	}
	port = int(ee)
	ip = net.IPv4(byte(aa), byte(bb), byte(cc), byte(dd))
	return
}