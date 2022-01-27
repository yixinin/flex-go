package topic

import (
	"context"
	"sync"

	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
)

type ReturnSender struct {
	locker     sync.RWMutex
	publishers map[string]*pubsub.Publisher
}

func (m *ReturnSender) Send(ctx context.Context, msg message.Message) (err error) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	if pub, ok := m.publishers[msg.PeerId()]; ok {
		return pub.Send(ctx, msg)
	}
	return nil
}

func (m *ReturnSender) OnSubJoin(ctx context.Context, sub *pubsub.Subscriber) {}
func (m *ReturnSender) OnSubLeave(ctx context.Context, id string)             {}
func (m *ReturnSender) OnPubJoin(ctx context.Context, pub *pubsub.Publisher) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.publishers[pub.Id()] = pub
}
func (m *ReturnSender) OnPubLeave(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.publishers, id)
}
