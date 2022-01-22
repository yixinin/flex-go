package topic

import (
	"context"

	"github.com/yixinin/flex/client"
)

func (m *TopicManager) AddPub(ctx context.Context, pub *client.Publisher) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.publishers[pub.Id()] = pub
	m.router.OnPubJoin(ctx, pub)
}
func (m *TopicManager) DelPub(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.publishers, id)
	m.router.OnPubLeave(ctx, id)
}
func (m *TopicManager) GetPub(key string) (*client.Publisher, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	pub, ok := m.publishers[key]
	return pub, ok
}

func (m *TopicManager) ForeachPub(ctx context.Context, f func(id string, pub *client.Publisher)) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	for k, v := range m.publishers {
		f(k, v)
	}
}

func (m *TopicManager) AddSub(ctx context.Context, sub *client.Subscriber) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.subscribers[sub.Id()] = sub
	m.router.OnSubJoin(ctx, sub)
}
func (m *TopicManager) DelSub(ctx context.Context, id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.subscribers, id)
	m.router.OnSubLeave(ctx, id)
}
func (m *TopicManager) GetSub(key string) (*client.Subscriber, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	sub, ok := m.subscribers[key]
	return sub, ok
}
func (m *TopicManager) ForeachSub(ctx context.Context, f func(id string, sub *client.Subscriber)) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	for k, v := range m.subscribers {
		f(k, v)
	}
}
