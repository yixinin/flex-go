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
	wg     sync.WaitGroup
	topics map[string]*topic.TopicManager

	beforeCtx      context.Context
	beforeShutdown func()

	ctx      context.Context
	shutdown func()

	afterCtx      context.Context
	afterShutdown func()
}

func NewManager(rawCtx context.Context, tps []topic.Config) *Manager {
	var m = &Manager{
		topics: make(map[string]*topic.TopicManager, len(tps)),
	}
	m.init(rawCtx, tps)
	return m
}

func (m *Manager) BeforeShutdown() {
	if m.beforeShutdown != nil {
		m.beforeShutdown()
	}
}

func (m *Manager) ShutDown() {
	if m.shutdown != nil {
		m.shutdown()
	}
}

func (m *Manager) AfterShutdown() {
	if m.afterShutdown != nil {
		m.afterShutdown()
	}
}

// exit when before shutdown called, stop accept
func (m *Manager) Run(rawCtx context.Context) error {
	err := m.listen()
	m.wg.Wait()
	return err
}

func (m *Manager) init(rawCtx context.Context, tps []topic.Config) {
	m.beforeCtx, m.beforeShutdown = context.WithCancel(rawCtx)
	m.ctx, m.shutdown = context.WithCancel(rawCtx)
	m.afterCtx, m.afterShutdown = context.WithCancel(rawCtx)

	for _, t := range tps {
		tp := topic.NewTopicManager(t.RouterName, t.BufferName)
		m.topics[t.Topic] = tp
		m.wg.Add(1)
		go func(tp *topic.TopicManager) {
			defer m.wg.Done()
			_ = tp.Run(m.ctx)
		}(tp)
	}
}

func (m *Manager) listen() error {
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 4569})
	if err != nil {
		return err
	}

	for {
		select {
		case <-m.beforeCtx.Done():
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
				logger.Errorf(m.beforeCtx, "read conn header error:%v", err)
				continue
			}
			if n != 1 {
				// ignore conn
				logger.Warnf(m.beforeCtx, "connection message:%s size:%d not match", headerBuf, n)
				continue
			}
			var buf = make([]byte, headerBuf[0])
			connMessage, ok := message.UnmarshalConnMessage(buf)
			if !ok {
				logger.Warnf(m.beforeCtx, "unknown connection message:%s", buf)
				continue
			}
			tm, ok := m.topics[connMessage.Topic]
			if !ok {
				logger.Warnf(m.beforeCtx, "no such topic:%s", connMessage.Topic)
				continue
			}

			switch connMessage.Type {
			case message.ClientTypeSub:
				// exit when shutdown called, stop recv new message, app exit
				ctx, cancel := context.WithCancel(m.ctx)
				sub := pubsub.NewSubscriber(conn, connMessage, cancel)
				sub.Recv(ctx, tm.Channel())
				tm.AddSub(ctx, sub)
			case message.ClientTypePub:
				ctx, cancel := context.WithCancel(m.ctx)
				pub := pubsub.NewPublisher(conn, connMessage, cancel)
				pub.Recv(ctx, tm.Channel())
				tm.AddPub(ctx, pub)
			}
		}
	}
}
