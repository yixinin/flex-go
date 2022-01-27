package client

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
)

type ConnManager interface {
	OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr)
	OnAddrLeave(ctx context.Context, id string)
}

type Conn struct {
	locker           sync.RWMutex
	Mode             string
	topic            string
	clientType       message.ClientType
	servers          map[string]*Server
	waitCloseServers map[string]*Server
	ch               chan message.Message
}

func (c *Conn) AddServer(s *Server) {
	c.locker.Lock()
	defer c.locker.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go s.recv(ctx, c.ch)
	go s.heartBeat(ctx)
	c.servers[s.id] = s

}

func (c *Conn) DelServer(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	s, ok := c.servers[id]
	if ok && s != nil && s.cancel != nil {
		s.Close(context.Background())
		c.waitCloseServers[id] = s
	}
	delete(c.servers, id)
}

func (c *Conn) DropServer(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	s, ok := c.waitCloseServers[id]
	if ok && s != nil {
		s.Drop()
	}
	delete(c.waitCloseServers, id)
}

func (c *Conn) GetServer(id string) (*Server, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	s, ok := c.servers[id]
	return s, ok
}
func (c *Conn) GetWaitCloseServer(id string) (*Server, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	s, ok := c.waitCloseServers[id]
	return s, ok
}

func (c *Conn) ForeachServers(f func(string, *Server)) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for k, v := range c.servers {
		f(k, v)
	}
}
func (c *Conn) OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr) {
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
	c.AddServer(server)

	<-ctx.Done()
}

func (c *Conn) OnAddrLeave(ctx context.Context, id string) {
	c.DelServer(id)
}

func (c *Conn) Ack(ctx context.Context, id, key, groupKey string) {
	s, ok := c.GetServer(id)
	if !ok {
		s, ok = c.GetWaitCloseServer(id)
	}
	if ok {
		err := s.Send(ctx, message.NewAckMessage(key, groupKey))
		if err != nil {
			logger.Error(ctx, err)
		}
	}
}
