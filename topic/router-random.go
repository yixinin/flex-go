package topic

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RandomSender struct {
	locker      sync.RWMutex
	subKeys     []string
	subscribers map[string]*client.Subscriber
}

func (m *RandomSender) syncKeys() {
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

func (m *RandomSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	m.syncKeys()
	key := m.subKeys[random()%len(m.subKeys)]
	if sub, ok := m.subscribers[key]; ok {
		err = sub.Send(ctx, msg)
	}
	return
}

func (m *RandomSender) OnSubJoin(ctx context.Context, sub *client.Subscriber) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.subscribers[sub.Id()] = sub
}
func (m *RandomSender) OnSubLeave(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.subscribers, id)
}
func (m *RandomSender) OnPubJoin(ctx context.Context, pub *client.Publisher) {}
func (m *RandomSender) OnPubLeave(ctx context.Context, id string)            {}

func random() int {
	return rand.Int()
}
