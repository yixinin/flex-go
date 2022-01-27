package message

type HeartBeat struct {
	peerId string
}

func NewHearbeatMessage() *HeartBeat {
	return &HeartBeat{}
}

func (m *HeartBeat) Id() string {
	return ""
}
func (m *HeartBeat) Group() string {
	return ""
}
func (m *HeartBeat) RawData() []byte {
	return nil
}

func (m *HeartBeat) PeerId() string {
	return m.peerId
}

func (m *HeartBeat) Marshal() []byte {
	var buf = make([]byte, HEADER_SIZE)
	buf[0] = byte(MessageTypeHeartBeat)
	return buf
}
