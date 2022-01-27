package message

import (
	"encoding/binary"
)

const (
	MessageTypeHeartBeat MessageType = 1
	MessageTypeClose     MessageType = 2
	MessageTypeConn      MessageType = 3
	MessageTypeRaw       MessageType = 4
	MessageTypeAck       MessageType = 5
)
const HEADER_SIZE = 9

type MessageType byte

type Header struct {
	Size        int
	MessageType MessageType
	peerId      string
}

func ParseHeader(peerId string, buf [HEADER_SIZE]byte) Header {
	msgType := MessageType(buf[0])
	switch msgType {
	case MessageTypeHeartBeat, MessageTypeClose, MessageTypeConn:
		return Header{
			peerId:      peerId,
			MessageType: msgType,
		}
	default:
		size := binary.BigEndian.Uint64(buf[1:])
		return Header{
			peerId:      peerId,
			Size:        int(size),
			MessageType: msgType,
		}
	}
}
