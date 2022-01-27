package message

type ClientType uint8

const (
	ClientTypePub = 1
	ClientTypeSub = 2
)

var ValidClientTypes = map[ClientType]struct{}{
	ClientTypePub: struct{}{},
	ClientTypeSub: struct{}{},
}

type ConnMessage struct {
	Topic string
	Type  ClientType
}

func UnmarshalConnMessage(buf []byte) (ConnMessage, bool) {
	var msg = ConnMessage{
		Type:  ClientType(buf[0]),
		Topic: string(buf[1:]),
	}
	if _, ok := ValidClientTypes[msg.Type]; !ok {
		return msg, false
	}

	return msg, false
}

func (msg ConnMessage) Marshal() []byte {
	var buf = make([]byte, 0, 1+len(msg.Topic)+1)
	buf[0] = byte(len(buf) - 1)
	buf[1] = byte(msg.Type)
	copy(buf[2:], []byte(msg.Topic))
	return buf
}
