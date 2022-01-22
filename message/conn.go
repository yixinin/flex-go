package message

type ClientType uint8

const (
	TypePub = 1
	TypeSub = 2
)

type ConnMessage struct {
	Topic string
	Type  ClientType
}

func ParseConnMessage(buf []byte) *ConnMessage {
	return &ConnMessage{}
}
