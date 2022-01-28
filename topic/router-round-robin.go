package topic

import (
	"context"
	"sync"

	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
)

type RoundRobinSender struct {
	locker      sync.RWMutex
	round       int
	subKeys     []string
	subscribers map[string]*pubsub.Subscriber
}

func (m *RoundRobinSender) syncKeys() {
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

func NewRoundRobinRouter() MessageRouter {
	return &RoundRobinSender{
		subscribers: make(map[string]*pubsub.Subscriber),
	}
}

func (m *RoundRobinSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	m.syncKeys()
	key := m.subKeys[m.round%len(m.subKeys)]
	if sub, ok := m.subscribers[key]; ok {
		err = sub.Send(ctx, msg)
	}
	m.round++
	return
}

func (m *RoundRobinSender) OnSubJoin(ctx context.Context, sub *pubsub.Subscriber) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.subscribers[sub.Id()] = sub
}
func (m *RoundRobinSender) OnSubLeave(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.subscribers, id)
}

func (m *RoundRobinSender) OnPubJoin(ctx context.Context, pub *pubsub.Publisher) {}
func (m *RoundRobinSender) OnPubLeave(ctx context.Context, id string)            {}
