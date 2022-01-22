package main

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
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
	var buf = make([]byte, 1024)

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
			n, err := conn.Read(buf[:])
			if err != nil {
				logger.Errorf(ctx, "read conn header error:%v", err)
			}
			if n != message.HEADER_SIZE || err != nil {
				// ignore conn
				logger.Warnf(ctx, "connection message:%s error:%v or size:%d not match", buf[:], err, n)
				continue
			}
			connMessage := message.ParseConnMessage(buf)
			tm, ok := m.topics[connMessage.Topic]
			if !ok {
				logger.Warnf(ctx, "no such topic:%s error:%v", connMessage.Topic)
				continue
			}

			switch connMessage.Type {
			case message.TypeSub:
				ctx, cancel := context.WithCancel(m.delayCtx)
				sub := client.NewSubscriber(conn, connMessage, cancel)
				sub.Recv(ctx, tm.Channel())
				tm.AddSub(m.delayCtx, sub)
			case message.TypePub:
				ctx, cancel := context.WithCancel(m.delayCtx)
				pub := client.NewPublisher(conn, connMessage, cancel)
				pub.Recv(ctx, tm.Channel())
				tm.AddPub(m.delayCtx, pub)
			}
		}
	}
}
