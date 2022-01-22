package topic

import "context"

func (m *TopicManager) AddGroup(ctx context.Context, key string, g *Group) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.groups[key] = g
	go g.Run(ctx, m.sendCh)
}

func (m *TopicManager) DelGroup(ctx context.Context, key string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	g, ok := m.groups[key]
	if ok {
		g.Close()
	}
	delete(m.groups, key)
}

func (m *TopicManager) GetGroup(key string) (*Group, bool) {
	g, ok := m.groups[key]
	return g, ok
}

func (m *TopicManager) ForeachGroup(ctx context.Context, f func(id string, group *Group)) {
	m.locker.RLock()
	defer m.locker.RUnlock()
	for k, v := range m.groups {
		f(k, v)
	}
}
