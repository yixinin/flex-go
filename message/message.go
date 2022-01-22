package message

type MessageStatus uint8

const (
	StatusNone    = 0
	StatusWaitAck = 1
	StatusAcked   = 2
)

type Message interface {
	Id() string
	ClientId() string
	Group() string
	RawData() []byte
	Status() MessageStatus
	SetStatus(MessageStatus)
}

type MergeIface interface {
	Merge(msg Message) bool
}
type SplitIface interface {
	Split() []Message
}

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
