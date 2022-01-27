package message

import "encoding/binary"

type AckMessage struct {
	Key      string `json:"key"`
	GroupKey string `json:"group_key"`
	peerId   string
}

func NewAckMessage(key, groupKey string) Message {
	return &AckMessage{
		Key:      key,
		GroupKey: groupKey,
	}
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

func (m *AckMessage) PeerId() string {
	return m.peerId
}

func (m *AckMessage) Marshal() []byte {
	var size = len(m.Key) + len(m.GroupKey)
	var buf = make([]byte, HEADER_SIZE+size)
	buf[0] = byte(MessageTypeAck)
	binary.BigEndian.PutUint64(buf[1:HEADER_SIZE], uint64(size))
	copy(buf[HEADER_SIZE:HEADER_SIZE+len(m.Key)], []byte(m.Key))
	copy(buf[HEADER_SIZE+len(m.Key):], []byte(m.GroupKey))
	return buf
}
