package message

const (
	TypeHeartBeat MessageType = 1
	HEADER_SIZE               = 8
)

type MessageType byte

type Header struct {
	Size        int
	MessageType MessageType
}

func ParseHeader(buf [HEADER_SIZE]byte) Header {
	return Header{}
}
