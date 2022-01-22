package message

type AckMessage struct {
	Key      string
	GroupKey string
	clientId string
}

func (m *AckMessage) Id() string {
	return m.Key
}
func (m *AckMessage) Group() string {
	return m.GroupKey
}
func (m *AckMessage) RawData() []byte {
	return nil
}

func (m *AckMessage) ClientId() string {
	return m.clientId
}

func (m *AckMessage) Status() MessageStatus {
	return StatusNone
}

func (m *AckMessage) SetStatus(status MessageStatus) {

}
