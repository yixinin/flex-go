package message

type HeartBeat struct {
	clientId string
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

func (m *HeartBeat) Status() MessageStatus {
	return StatusNone
}

func (m *HeartBeat) SetStatus(status MessageStatus) {

}
