package conn

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/ttl"
)

type Server struct {
	id     string
	conn   net.Conn
	ttl    int64
	cancel func()
}

func (s *Server) Send(ctx context.Context, msg message.Message) error {
	_, err := s.conn.Write(msg.Marshal())
	return err
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
			switch header.MessageType {
			case message.MessageTypeHeartBeat:
				s.ttl = ttl.NextTTL()
			case message.MessageTypeClose:
				s.Close(ctx)
			case message.MessageTypeRaw:
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

func (s *Server) TTL() int64 {
	return s.ttl
}

// call close when application recv exit sig
func (s *Server) Close(ctx context.Context) {
	if err := s.Send(ctx, message.NewCloseMessage()); err != nil {
		logger.Error(ctx, err)
	}
}

// beteewn close and drop, client is allowed to send ack message to server
// call drop whenn application ready to exit
func (s *Server) Drop() {
	if s.cancel != nil {
		s.cancel()
	}
	if s.conn != nil {
		s.conn.Close()
	}
}
