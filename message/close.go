package message

type CloseMessage struct {
	clientId string
}

func NewCloseMessage() Message {
	return &CloseMessage{}
}

func (m *CloseMessage) Id() string {
	return ""
}
func (m *CloseMessage) Group() string {
	return ""
}
func (m *CloseMessage) RawData() []byte {
	return nil
}

func (m *CloseMessage) ClientId() string {
	return m.clientId
}

func (m *CloseMessage) Marshal() []byte {
	var buf = make([]byte, 1)
	buf[0] = byte(TypeClose)
	return buf
}
