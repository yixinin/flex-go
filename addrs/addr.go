package addrs

import (
	"context"
	"encoding/binary"
	"net"
	"strings"

	"github.com/yixinin/flex/logger"
)

func Marshal(addr *net.TCPAddr) []byte {
	var data = make([]byte, 6)
	ip := addr.IP.To4()
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(addr.Port))
	copy(data[:4], ip)
	copy(data[4:], port)
	return data
}

func Unmarshal(data []byte) *net.TCPAddr {
	var port = binary.BigEndian.Uint16(data[4:])
	var addr = &net.TCPAddr{
		IP:   data[:4],
		Port: int(port),
	}
	return addr
}

func ParseKv(ctx context.Context, key, val []byte) (id string, addr *net.TCPAddr) {
	keys := strings.Split(string(key), "/")
	if len(keys) != 2 || len(val) == 0 {
		logger.Warnf(ctx, "unknown key:%s val:%s", key, val)
		return
	}
	id = keys[2]
	addr = Unmarshal(val)
	return
}
