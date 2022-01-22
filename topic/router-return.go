package topic

import (
	"context"
	"sync"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type ReturnSender struct {
	locker     sync.RWMutex
	publishers map[string]*client.Publisher
}

func (m *ReturnSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	if pub, ok := m.publishers[msg.ClientId()]; ok {
		return pub.Send(ctx, msg)
	}
	return nil
}

func (m *ReturnSender) OnSubJoin(ctx context.Context, sub *client.Subscriber) {}
func (m *ReturnSender) OnSubLeave(ctx context.Context, id string)             {}
func (m *ReturnSender) OnPubJoin(ctx context.Context, pub *client.Publisher) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.publishers[pub.Id()] = pub
}
func (m *ReturnSender) OnPubLeave(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.publishers, id)
}
