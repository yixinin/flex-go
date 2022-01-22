package topic

import (
	"context"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type RoundRobinSender struct {
	round       int
	subKeys     []string
	subscribers map[string]*client.Subscriber
}

func (m RoundRobinSender) syncKeys() {
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

func (m *RoundRobinSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.syncKeys()
	key := m.subKeys[m.round%len(m.subKeys)]
	if sub, ok := m.subscribers[key]; ok {
		err = sub.Send(ctx, msg)
	}
	m.round++
	return
}
