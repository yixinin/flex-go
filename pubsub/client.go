package pubsub

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"github.com/yixinin/flex/logger"
	"github.com/yixinin/flex/message"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	conn         net.Conn
	disconnected bool
	topic        string
	id           string
	ttl          int64
	cancel       func()
}

func newClient(conn net.Conn, connMessage *message.ConnMessage, cancel func()) *Client {
	return &Client{
		conn:   conn,
		topic:  connMessage.Topic,
		id:     primitive.NewObjectID().Hex(),
		cancel: cancel,
	}
}

func (c *Client) Id() string {
	return c.id
}

// unix nano
func (c *Client) TTL() int64 {
	return c.ttl
}

func (c *Client) Close() {
	c.cancel()
}

func (c *Client) Disconnected() bool {
	return c.disconnected
}

func (c *Client) Send(ctx context.Context, msg message.Message) error {
	_, err := c.conn.Write(msg.RawData())
	return err
}

func (c *Client) Recv(ctx context.Context, ch chan message.Message) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error(ctx, r)
			}
			c.disconnected = true
			c.conn.Close()
		}()

		var headBuf [message.HEADER_SIZE]byte

		for {
			select {
			case <-ctx.Done():
				logger.Infof(ctx, "client: %+v disconnected", c)
				return
			default:
				c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
				n, err := c.conn.Read(headBuf[:])
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				if err != nil || n != message.HEADER_SIZE {
					logger.Errorf(ctx, "client: %+v recv error:%v", c, err)
					return
				}
				var header = message.ParseHeader(headBuf)
				switch header.MessageType {
				case message.TypeHeartBeat:
					c.ttl = time.Now().Add(time.Second).UnixNano()
				default:
					var buf = make([]byte, header.Size)
					n, err = c.conn.Read(buf)
					if err != nil || n != header.Size {
						logger.Errorf(ctx, "client: %+v recv error:%v", c, err)
						return
					}
					var msg = message.ToRawMessage(headBuf[:], buf)
					ch <- msg
				}
			}
		}

	}()
}
