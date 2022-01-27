package topic

import (
	"context"
	"sync"
	"time"

	"github.com/yixinin/flex/buffers"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
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
	wg     sync.WaitGroup

	router      MessageRouter
	subscribers map[string]*pubsub.Subscriber
	publishers  map[string]*pubsub.Publisher
	groups      map[string]*Group

	recvCh chan message.Message
	sendCh chan message.Message

	newBuffer func() buffers.Buffer
}

func NewTopicManager(routerName, bufferName string) *TopicManager {
	return &TopicManager{
		router:      NewRouter(routerName),
		subscribers: make(map[string]*pubsub.Subscriber, 1),
		publishers:  make(map[string]*pubsub.Publisher, 1),
		groups:      make(map[string]*Group, 1),
		recvCh:      make(chan message.Message),
		sendCh:      make(chan message.Message),
		newBuffer:   buffers.NewBufferFunc(bufferName),
	}
}

func (m *TopicManager) Channel() chan message.Message {
	return m.recvCh
}

func (m *TopicManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf(ctx, "topic manager run with error:%v", r)
		}
		wg.Done()
	}()

	m.wg.Add(1)
	go m.recv(ctx)

	m.wg.Add(1)
	go m.send(ctx)

	go m.checkTTL(ctx)
	go m.checkConn(ctx)

	m.wg.Wait()
}

func (m *TopicManager) recv(ctx context.Context) {
	defer m.wg.Done()
	for msg := range m.recvCh {
		execFunc(ctx, func() {
			groupKey := msg.Group()
			group, ok := m.GetGroup(groupKey)
			if !ok {
				group = newGroup(groupKey, m.newBuffer)
				m.wg.Add(1)
				m.AddGroup(ctx, groupKey, group)
			}
			group.TTL = time.Now().Unix() + TTL
			group.ch <- msg
		})
	}
}

func (m *TopicManager) send(ctx context.Context) {
	defer m.wg.Done()
	for msg := range m.sendCh {
		execFunc(ctx, func() {
			if err := m.router.Send(ctx, msg); err != nil {
				logger.Errorf(ctx, "send message error:%", err)
			}
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

				m.ForeachGroup(ctx, func(id string, group *Group) {
					if group.TTL <= nowUnix {
						waitDeleteGroups = append(waitDeleteGroups, id)
					}
				})

				for _, k := range waitDeleteGroups {
					m.DelGroup(ctx, k)
				}

				time.Sleep(100 * time.Millisecond)
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

				m.ForeachSub(func(id string, sub *pubsub.Subscriber) {
					if sub.TTL() < nowUnix {
						sub.Drop()
						waitDeleteSubIds = append(waitDeleteSubIds, id)
					}
					if sub.Disconnected() {
						waitDeleteSubIds = append(waitDeleteSubIds, sub.Id())
					}
				})
				m.ForeachPub(func(id string, pub *pubsub.Publisher) {
					if pub.TTL() < nowUnix {
						pub.Drop()
						waitDeletePubIds = append(waitDeletePubIds, id)
					}
					if pub.Disconnected() {
						waitDeletePubIds = append(waitDeletePubIds, pub.Id())
					}
				})

				for _, k := range waitDeletePubIds {
					m.DelPub(ctx, k)
				}
				for _, k := range waitDeleteSubIds {
					m.DelSub(ctx, k)
				}
			})
		}
	}
}
