package message

import "encoding/binary"

type MessageStatus uint8

type Message interface {
	Id() string
	ClientId() string
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

func Unmarshal(header Header, buf []byte, clinetId string) (Message, error) {
	switch header.MessageType {
	case TypeHeartBeat:
		return &HeartBeat{
			clientId: clinetId,
		}, nil
	case TypeRaw:
		return ToRawMessage(header, buf), nil
	}
	return nil, nil
}

func Marshal(bufs [][]byte) []byte {
	var size = uint64(0)
	for _, buf := range bufs {
		size += uint64(len(buf))
	}
	var buf = make([]byte, HEADER_SIZE+size)
	buf[0] = byte(TypeAck)
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
