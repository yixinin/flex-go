package message

import "encoding/binary"

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

func (m *AckMessage) Marshal() []byte {
	var size = len(m.Key) + len(m.GroupKey)
	var buf = make([]byte, HEADER_SIZE+size)
	buf[0] = byte(TypeAck)
	var sizeBuf = make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBuf, uint64(size))
	copy(buf[1:HEADER_SIZE], sizeBuf)
	copy(buf[HEADER_SIZE:HEADER_SIZE+len(m.Key)], []byte(m.Key))
	copy(buf[HEADER_SIZE+len(m.Key):], []byte(m.GroupKey))
	return buf
}
