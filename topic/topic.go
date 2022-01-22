package topic

import (
	"context"
	"sync"
	"time"

	"github.com/yixinin/flex/buffers"
	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
)

var execFunc = func(ctx context.Context, f func()) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(ctx, "run func error:%v", r)
		}
	}()

	f()
}

var TTL = int64(10 * 60)

type TopicManager struct {
	locker sync.RWMutex

	router      MessageRouter
	subscribers map[string]*client.Subscriber
	publishers  map[string]*client.Publisher

	recvCh     chan message.Message
	sendCh     chan message.Message
	groups     map[string]*Group
	newBuffer  func() buffers.Buffer
	Distribute message.DistributeType
}

func NewTopicManager() *TopicManager {
	return &TopicManager{}
}

func (m *TopicManager) AddGroup(ctx context.Context, key string, g *Group) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.groups[key] = g
	go g.Run(ctx, m.sendCh)
}

func (m *TopicManager) Channel() chan message.Message {
	return m.recvCh
}

func (m *TopicManager) AddPub(ctx context.Context, pub *client.Publisher) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.publishers[pub.Id()] = pub
}
func (m *TopicManager) DelPub(id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.publishers, id)
}
func (m *TopicManager) AddSub(ctx context.Context, sub *client.Subscriber) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.subscribers[sub.Id()] = sub
}
func (m *TopicManager) DelSub(id string) {
	m.locker.Lock()
	defer m.locker.Unlock()
	delete(m.subscribers, id)
}

func (m *TopicManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(ctx, "topic manager run with error:%v", r)
		}
		wg.Done()
	}()
	go m.recv(ctx)
	go m.checkTTL(ctx)
	go m.checkConn(ctx)

	go m.send(ctx)
}

func (m *TopicManager) send(ctx context.Context) {
	for msg := range m.sendCh {
		execFunc(ctx, func() {
			switch m.Distribute {
			case message.HttpRequest:
				if err := m.router.Send(ctx, msg); err != nil {
					logger.Errorf(ctx, "send message error:%", err)
				}
			default:
				if err := m.router.Send(ctx, msg); err != nil {
					logger.Errorf(ctx, "send message error:%", err)
				}
			}
		})
	}
}

func (m *TopicManager) recv(ctx context.Context) {
	for msg := range m.recvCh {
		execFunc(ctx, func() {
			groupKey := msg.Group()
			group, ok := m.groups[groupKey]
			if !ok {
				group = newGroup(groupKey, m.newBuffer)
				m.groups[groupKey] = group
			}
			group.TTL = time.Now().Unix() + TTL
			group.ch <- msg
		})
	}
}

func (m *TopicManager) checkTTL(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			execFunc(ctx, func() {
				var waitDeleteGroups = make([]string, 0)
				var nowUnix = time.Now().Unix()
				for key, group := range m.groups {
					if group.TTL <= nowUnix {
						waitDeleteGroups = append(waitDeleteGroups, key)
					}
				}
				for _, k := range waitDeleteGroups {
					delete(m.groups, k)
				}
				time.Sleep(1 * time.Second)
			})
		}
	}
}

func (m *TopicManager) checkConn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			execFunc(ctx, func() {
				var waitDeletePubIds = make([]string, 0)
				var waitDeleteSubIds = make([]string, 0)
				var nowUnix = time.Now().UnixNano()
				for id, sub := range m.subscribers {
					if sub.TTL() < nowUnix {
						sub.Close()
						waitDeleteSubIds = append(waitDeleteSubIds, id)
					}
				}
				for id, pub := range m.publishers {
					if pub.TTL() < nowUnix {
						pub.Close()
						waitDeletePubIds = append(waitDeletePubIds, id)
					}
				}
				for _, k := range waitDeletePubIds {
					delete(m.subscribers, k)
				}
				for _, k := range waitDeleteSubIds {
					delete(m.subscribers, k)
				}
			})
		}
	}
}
