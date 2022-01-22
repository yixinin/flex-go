package topic

import (
	"context"
	"hash/crc32"
	"sync"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type HashSender struct {
	locker      sync.RWMutex
	subKeys     []string
	subscribers map[string]*client.Subscriber
}

func NewHashRouter() MessageRouter {
	return &HashSender{
		subscribers: make(map[string]*client.Subscriber),
	}
}

func (m *HashSender) syncKeys() {
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
	m.locker.RLock()
	defer m.locker.RUnlock()
	m.syncKeys()
	key := m.subKeys[hash(msg.Group())%len(m.subKeys)]
	if sub, ok := m.subscribers[key]; ok {
		err = sub.Send(ctx, msg)
	}
	return
}

func (m *HashSender) OnSubJoin(ctx context.Context, sub *client.Subscriber) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.subscribers[sub.Id()] = sub
}
func (m *HashSender) OnSubLeave(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.subscribers, id)
}

func (m *HashSender) OnPubJoin(ctx context.Context, pub *client.Publisher) {}
func (m *HashSender) OnPubLeave(ctx context.Context, id string)            {}

func hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}
