package client

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
)

type ConnManager interface {
	OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr)
	OnAddrLeave(ctx context.Context, id string)
}

type Conn struct {
	Mode       string
	topic      string
	clientType message.ClientType
	servers    map[string]*Server
	ch         chan message.Message
}

type Server struct {
	id   string
	conn net.Conn
	TTL  int64
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
		TTL:  time.Now().Add(time.Second).UnixNano(),
	}
	c.servers[id] = server

	go c.recv(ctx, server)
}

func (c *Conn) recv(ctx context.Context, s *Server) {
	var headerBuf [message.HEADER_SIZE]byte
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.conn.SetDeadline(time.Now().Add(2 * time.Second))
			n, err := s.conn.Read(headerBuf[:])
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}
			if err != nil {
				return
			}
			if n != message.HEADER_SIZE {
				continue
			}
			header := message.ParseHeader(headerBuf)
			if header.MessageType == message.TypeHeartBeat {
				s.TTL = time.Now().Add(time.Second).UnixNano()
				continue
			}
			var buf = make([]byte, header.Size)
			_, err = s.conn.Read(buf)
			if err != nil {
				return
			}
			msg, err := message.Unmarshal(header, buf, s.id)
			if err != nil {
				continue
			}
			c.ch <- msg
		}
	}
}

func (c *Conn) heartBeat(ctx context.Context, conn net.Conn) {
	var tick = time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			_, err := conn.Write(message.NewHearbeatMessage().Marshal())
			if err != nil {
				logger.Errorf(ctx, "heartbeat error:%v", err)
				return
			}
		}
	}
}

func (s *Server) SendClose(ctx context.Context) {
	s.conn.Write(message.NewCloseMessage().Marshal())
}
