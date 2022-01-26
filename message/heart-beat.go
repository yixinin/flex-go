package message

type HeartBeat struct {
	clientId string
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

func (m *HeartBeat) Distribute() DistributeType {
	return DisNone
}

func (m *HeartBeat) ClientId() string {
	return m.clientId
}

func (m *HeartBeat) Marshal() []byte {
	var buf = make([]byte, 1)
	buf[0] = byte(TypeHeartBeat)
	return buf
}
