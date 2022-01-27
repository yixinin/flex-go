package message

type CloseMessage struct {
	peerId string
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

func (m *CloseMessage) PeerId() string {
	return m.peerId
}

func (m *CloseMessage) Marshal() []byte {
	var buf = make([]byte, HEADER_SIZE)
	buf[0] = byte(MessageTypeClose)
	return buf
}
