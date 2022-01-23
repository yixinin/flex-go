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

func Unmarshal(header Header, buf []byte, clinetId string) (Message, error) {
	switch header.MessageType {
	case TypeHeartBeat:
		return &HeartBeat{
			clientId: clinetId,
		}, nil
	case TypeRaw:
		return ParseConnMessage(buf), nil
	}
	return nil, nil
}
