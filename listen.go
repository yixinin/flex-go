package main

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
	"github.com/yixinin/flex/topic"
)

type Manager struct {
	locker   sync.RWMutex
	wg       *sync.WaitGroup
	delayCtx context.Context
	topics   map[string]*topic.TopicManager
}

func NewManager(delayCtx context.Context) *Manager {
	return &Manager{
		delayCtx: delayCtx,
		topics:   make(map[string]*topic.TopicManager),
		wg:       &sync.WaitGroup{},
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

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			lis.SetDeadline(time.Now().Add(time.Second))
			conn, err := lis.Accept()
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}
			if err != nil {
				return err
			}
			var headerBuf = make([]byte, 1)
			n, err := conn.Read(headerBuf)
			if err != nil {
				logger.Errorf(ctx, "read conn header error:%v", err)
				continue
			}
			if n != 1 {
				// ignore conn
				logger.Warnf(ctx, "connection message:%s size:%d not match", headerBuf, n)
				continue
			}
			var buf = make([]byte, headerBuf[0])
			connMessage, ok := message.UnmarshalConnMessage(buf)
			if !ok {
				logger.Warnf(ctx, "unknown connection message:%s", buf)
				continue
			}
			tm, ok := m.topics[connMessage.Topic]
			if !ok {
				logger.Warnf(ctx, "no such topic:%s", connMessage.Topic)
				continue
			}

			switch connMessage.Type {
			case message.ClientTypeSub:
				ctx, cancel := context.WithCancel(m.delayCtx)
				sub := pubsub.NewSubscriber(conn, connMessage, cancel)
				sub.Recv(ctx, tm.Channel())
				tm.AddSub(m.delayCtx, sub)
			case message.ClientTypePub:
				ctx, cancel := context.WithCancel(m.delayCtx)
				pub := pubsub.NewPublisher(conn, connMessage, cancel)
				pub.Recv(ctx, tm.Channel())
				tm.AddPub(m.delayCtx, pub)
			}
		}
	}
}
