package message

import (
	"encoding/binary"
)

const (
	TypeHeartBeat MessageType = 1
	TypeClose     MessageType = 2
	TypeConn      MessageType = 3
	TypeRaw       MessageType = 4
)
const HEADER_SIZE = 9

type MessageType byte

type Header struct {
	Size        int
	MessageType MessageType
}

func ParseHeader(buf [HEADER_SIZE]byte) Header {
	size := binary.BigEndian.Uint64(buf[:HEADER_SIZE-1])
	return Header{
		Size:        int(size),
		MessageType: MessageType(buf[HEADER_SIZE-1]),
	}
}
