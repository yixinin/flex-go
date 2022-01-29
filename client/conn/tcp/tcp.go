package tcp

import (
	"context"
	"errors"
	"hash/crc32"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yixinin/flex/client/config"
	"github.com/yixinin/flex/client/conn"
	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/ttl"
)

type Config struct {
	Topic  string
	Pubsub string
}

func (c Config) Check() bool {
	if c.Pubsub == "" || c.Topic == "" {
		return false
	}
	return true
}

type TcpConnManager struct {
	locker           sync.RWMutex
	topic            string
	clientType       message.ClientType
	servers          map[string]*Server
	waitCloseServers map[string]*Server
	ch               chan message.Message
}

func NewTcpConnManager(conf config.Config) conn.ConnManager {
	c, _ := conf.(Config)
	return &TcpConnManager{
		topic: c.Topic,
		clientType: func(pubsub string) message.ClientType {
			switch pubsub {
			case "pub":
				return message.ClientTypePub
			case "sub":
				return message.ClientTypeSub
			}
			panic("unknown client type:" + pubsub)
		}(c.Pubsub),
		servers:          make(map[string]*Server),
		waitCloseServers: make(map[string]*Server),
		ch:               make(chan message.Message, 1024),
	}
}

func (c *TcpConnManager) addServer(ctx context.Context, s *Server) {
	c.locker.Lock()
	defer c.locker.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go s.recv(ctx, c.ch)
	go s.heartBeat(ctx)
	c.servers[s.id] = s
}

func (c *TcpConnManager) delServer(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	s, ok := c.servers[id]
	ws, wok := c.waitCloseServers[id]
	if ok && s != nil {
		c.waitCloseServers[id] = s
	}

	if wok && ws != nil {
		ws.Drop()
		delete(c.waitCloseServers, id)
	}
	delete(c.servers, id)
}

func (c *TcpConnManager) DropAllServer(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	s, ok := c.waitCloseServers[id]
	if ok && s != nil {
		s.Drop()
	}
	delete(c.waitCloseServers, id)

	s, ok = c.servers[id]
	if ok && s != nil {
		s.Drop()
	}
	delete(c.servers, id)
}

func (c *TcpConnManager) GetServer(id string) (*Server, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	s, ok := c.servers[id]
	return s, ok
}
func (c *TcpConnManager) GetWaitCloseServer(id string) (*Server, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	s, ok := c.waitCloseServers[id]
	return s, ok
}

func (c *TcpConnManager) ForeachServers(f func(string, *Server)) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for k, v := range c.servers {
		f(k, v)
	}
}

func (c *TcpConnManager) ForeachWaitCloseServers(f func(string, *Server)) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for k, v := range c.waitCloseServers {
		f(k, v)
	}
}

func (c *TcpConnManager) OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logger.Errorf(ctx, "connect to %s errro:%v", addr, err)
		return
	}

	data := message.ConnMessage{
		Topic: c.topic,
		Type:  c.clientType,
	}

	_, err = conn.Write(data.Marshal())
	if err != nil {
		return
	}

	var server = &Server{
		id:   id,
		conn: conn,
		ttl:  time.Now().Add(time.Second).UnixNano(),
	}
	c.addServer(ctx, server)
}

func (c *TcpConnManager) OnAddrLeave(ctx context.Context, id string) {
	c.delServer(id)
}

func (c *TcpConnManager) Send(ctx context.Context, msgid, groupKey string, payload []byte) error {
	var msg = message.NewRawMessage(msgid, groupKey, payload)
	var ids = make([]string, 0, len(c.servers))
	c.locker.RLock()
	defer c.locker.RUnlock()
	for k := range c.servers {
		ids = append(ids, k)
	}
	id := ids[hash(groupKey)%len(ids)]
	s, _ := c.GetServer(id)
	return s.Send(ctx, msg)
}
func (c *TcpConnManager) SendAsync(ctx context.Context, msgid, groupKey string, payload []byte) {
	go func() {
		err := c.Send(ctx, msgid, groupKey, payload)
		if err != nil {
			logger.Error(ctx, err)
		}
	}()
}

func (c *TcpConnManager) Recv(ctx context.Context, timeout time.Duration) (message.Message, error) {
	select {
	case <-time.After(timeout):
		return nil, os.ErrDeadlineExceeded
	case msg := <-c.ch:
		if msg == nil {
			return nil, errors.New("conn closed")
		}
		return msg, nil
	}
}

func (c *TcpConnManager) Ack(ctx context.Context, clientId, key, groupKey string) {
	s, ok := c.GetServer(clientId)
	if !ok {
		s, ok = c.GetWaitCloseServer(clientId)
	}
	if ok {
		err := s.Send(ctx, message.NewAckMessage(key, groupKey))
		if err != nil {
			logger.Error(ctx, err)
		}
	}
}

func hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

func (c *TcpConnManager) Run(ctx context.Context) error {
	c.checkTTL(ctx)
	return nil
}

func (c *TcpConnManager) checkTTL(ctx context.Context) {
	tick := ttl.NewTicker()
	defer tick.Stop()
	for {
		var now = time.Now().UnixMilli()
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			var waitCloseServers = make([]string, 0)
			var waitDropServers = make([]string, 0)
			c.ForeachServers(func(s1 string, s2 *Server) {
				if s2.ttl < now {
					waitCloseServers = append(waitCloseServers, s1)
				}
			})
			c.ForeachWaitCloseServers(func(s1 string, s2 *Server) {
				if s2.ttl < now {
					waitDropServers = append(waitDropServers, s1)
				}
			})
			for _, id := range waitCloseServers {
				c.DropAllServer(id)
			}
			for _, id := range waitDropServers {
				c.DropAllServer(id)
			}
		}
	}
}
