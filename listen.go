package main

import (
	"context"
	"net"
	"sync"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/topic"
)

type Manager struct {
	locker sync.RWMutex
	wg     *sync.WaitGroup
	topics map[string]*topic.TopicManager
}

func NewManager() *Manager {
	return &Manager{
		topics: make(map[string]*topic.TopicManager),
	}
}

func (m *Manager) Wait() {
	m.wg.Wait()
}

func (m *Manager) Run(ctx context.Context) {
	m.wg.Add(1)
	go m.Listen(ctx)
}

func (m *Manager) Listen(ctx context.Context) error {
	defer func() {
		m.wg.Done()
	}()
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 4569})
	if err != nil {
		return err
	}
	var buf = make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := lis.Accept()
			if err != nil {
				return err
			}
			n, err := conn.Read(buf[:])
			if err != nil {
				logger.Errorf(ctx, "read conn header error:%v", err)
			}
			if n != message.HEADER_SIZE || err != nil {
				// ignore conn
				continue
			}
			connMessage := message.ParseConnMessage(buf)
			tm, ok := m.topics[connMessage.Topic]
			if !ok {
				tm = topic.NewTopicManager()
				m.wg.Add(1)
				go tm.Run(ctx, m.wg)
				m.topics[connMessage.Topic] = tm
			}

			switch connMessage.Type {
			case message.TypeSub:
				ctx, cancel := context.WithCancel(ctx)
				sub := client.NewSubscriber(conn, connMessage, cancel)
				sub.Recv(ctx, nil)
				tm.AddSub(ctx, sub)
			case message.TypePub:
				ctx, cancel := context.WithCancel(ctx)
				pub := client.NewPublisher(conn, connMessage, cancel)
				pub.Recv(ctx, tm.Channel())
				tm.AddPub(ctx, pub)
			}
		}
	}
}
