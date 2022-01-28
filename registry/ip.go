package registry

import (
	"context"
	"net"

	"github.com/yixinin/flex/logger"
)

func GetLocalIP() net.IP {
	var ips = make([]net.IP, 0, 1)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Errorf(context.Background(), "get ip error:%v", err)
		return nil
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	for _, ip := range ips {
		return ip
	}
	return nil
}
