package topic

import (
	"context"
	"hash/crc32"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type HashSender struct {
	subKeys     []string
	subscribers map[string]*client.Subscriber
}

func (m HashSender) syncKeys() {
	var i int
	var keySize = len(m.subKeys)
	for k := range m.subscribers {
		if keySize > i {
			m.subKeys[i] = k
		} else {
			m.subKeys = append(m.subKeys, k)
		}
		i++
	}
}

func (m *HashSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.syncKeys()
	key := m.subKeys[hash(msg.Group())%len(m.subKeys)]
	if sub, ok := m.subscribers[key]; ok {
		err = sub.Send(ctx, msg)
	}
	return
}

func hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}
