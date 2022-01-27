package message

import (
	"encoding/binary"
	"encoding/json"
)

type RawMessage struct {
	Key      string `json:"key"`
	GroupKey string `json:"group_key"`
	Payload  []byte `json:"payload"`
	peerId   string
}

func NewRawMessage(key, groupKey string, payload []byte) Message {
	msg := &RawMessage{
		Key:      key,
		GroupKey: groupKey,
		Payload:  payload,
	}
	return msg
}

func newRawMessage(header Header, buf []byte) (Message, error) {
	msg := &RawMessage{
		peerId: header.peerId,
	}
	err := json.Unmarshal(buf, msg)
	return msg, err
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

func (m *RawMessage) PeerId() string {
	return m.peerId
}

func (m *RawMessage) Marshal() []byte {
	data, _ := json.Marshal(m)
	size := len(data)
	var buf = make([]byte, HEADER_SIZE+size)
	buf[0] = byte(MessageTypeRaw)
	binary.BigEndian.PutUint64(buf[1:HEADER_SIZE], uint64(size))
	copy(buf[HEADER_SIZE:], data)
	return buf
}
