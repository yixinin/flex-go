package message

type RawMessage struct {
	Key      string
	GroupKey string
	Payload  []byte
	clientId string
	status   MessageStatus
}

func ToRawMessage(headerBuf, buf []byte) *RawMessage {
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

func (m *RawMessage) Status() MessageStatus {
	return m.status
}

func (m *RawMessage) SetStatus(status MessageStatus) {
	m.status = status
}
