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

type Server struct {
	id     string
	conn   net.Conn
	TTL    int64
	cancel func()
}

func (s *Server) Send(ctx context.Context, msg message.Message) error {
	_, err := s.conn.Write(msg.Marshal())
	return err
}

func (s *Server) PreDisconnect(ctx context.Context) {
	err := s.Send(ctx, message.NewCloseMessage())
	if err != nil {
		logger.Error(ctx, err)
	}
}

func (s *Server) recv(ctx context.Context, ch chan message.Message) {
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
			header := message.ParseHeader(s.id, headerBuf)
			if header.MessageType == message.MessageTypeHeartBeat {
				s.TTL = time.Now().Add(time.Second).UnixNano()
				continue
			}
			var buf = make([]byte, header.Size)
			_, err = s.conn.Read(buf)
			if err != nil {
				logger.Error(ctx, err)
				return
			}
			msg, err := message.Unmarshal(header, buf)
			if err != nil {
				logger.Error(ctx, err)
				continue
			}
			ch <- msg
		}
	}
}

func (s *Server) heartBeat(ctx context.Context) {
	var tick = time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			_, err := s.conn.Write(message.NewHearbeatMessage().Marshal())
			if err != nil {
				logger.Errorf(ctx, "heartbeat error:%v", err)
				return
			}
		}
	}
}
