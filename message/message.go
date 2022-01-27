package message

import "encoding/binary"

type MessageStatus uint8

type Message interface {
	Id() string
	PeerId() string
	Group() string
	RawData() []byte
	Marshal() []byte
}

type MergeIface interface {
	Merge(msg Message) bool
}
type SplitIface interface {
	Split() []Message
}

func Unmarshal(header Header, buf []byte) (Message, error) {
	switch header.MessageType {
	case MessageTypeHeartBeat:
		return &HeartBeat{
			peerId: header.peerId,
		}, nil
	case MessageTypeRaw:
		return newRawMessage(header, buf)
	}
	return nil, nil
}

func Marshal(bufs [][]byte) []byte {
	var size = uint64(0)
	for _, buf := range bufs {
		size += uint64(len(buf))
	}
	var buf = make([]byte, HEADER_SIZE+size)
	buf[0] = byte(MessageTypeAck)
	var sizeBuf = make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBuf, uint64(size))
	copy(buf[1:HEADER_SIZE], sizeBuf)
	curIndex := HEADER_SIZE
	for _, v := range bufs {
		nextIndex := curIndex + len(v)
		copy(buf[curIndex:nextIndex], v)
		curIndex = nextIndex
	}
	return buf
}
