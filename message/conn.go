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

func (msg *ConnMessage) Marshal() []byte {
	var buf = make([]byte, 0, 1+len(msg.Topic)+1)
	buf[0] = byte(len(buf) - 1)
	copy(buf[1:len(msg.Topic)], []byte(msg.Topic))
	buf[len(msg.Topic)] = byte(msg.Type)
	return buf
}
