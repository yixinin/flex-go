package message

type RawMessage struct {
	Key      string
	GroupKey string
	Payload  []byte
	clientId string
}

func ToRawMessage(header Header, buf []byte) Message {
	return &RawMessage{}
}

func (m *RawMessage) Id() string {
	return m.Key
}
func (m *RawMessage) Group() string {
	return m.GroupKey
}
func (m *RawMessage) RawData() []byte {
	return m.Payload
}

func (m *RawMessage) ClientId() string {
	return m.clientId
}

func (m *RawMessage) Marshal() []byte {
	var buf = make([]byte, 0)
	return buf
}
