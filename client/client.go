package client

import (
	"context"
	"time"

	"github.com/yixinin/flex/client/event"
	"github.com/yixinin/flex/message"
)

type Client struct {
	connMgr ConnManager
	addrMgr AddrManager

	beforeCtx      context.Context
	beforeShutdown func()

	ctx      context.Context
	shutdown func()
}

func NewClient(app string) *Client {
	c := &Client{}
	return c
}

func (c *Client) SetConnManager(connMgr ConnManager) *Client {
	c.connMgr = connMgr
	return c
}

func (c *Client) SetAddrManager(addrMgr AddrManager, conf interface{}) *Client {
	c.addrMgr = addrMgr
	return c
}

func (c *Client) Run(rawCtx context.Context) error {
	if c.addrMgr == nil || c.connMgr == nil {
		panic("addr or conn manager is nil")
	}
	c.ctx, c.shutdown = context.WithCancel(rawCtx)
	c.beforeCtx, c.beforeShutdown = context.WithCancel(rawCtx)
	go c.addrMgr.Run(c.beforeCtx)
	go c.connMgr.Run(c.beforeCtx)
	c.onConnEvent(c.ctx)
	return nil
}

func (c *Client) onConnEvent(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-c.addrMgr.Event():
			switch ev.EventType {
			case event.EventAdd:
				c.connMgr.OnAddrJoin(ctx, ev.Id, ev.Addr)
			case event.EventDel:
				c.connMgr.OnAddrLeave(ctx, ev.Id)
			}
		}
	}
}

func (c *Client) BeforeShutdown() {
	c.beforeShutdown()
}
func (c *Client) Shutdown() {
	c.shutdown()
}

func (c *Client) AfterShutdown() {}

func (c *Client) Publish(ctx context.Context, msgid, groupKey string, payload []byte) error {
	return c.connMgr.Send(ctx, msgid, groupKey, payload)
}

func (c *Client) Recv(ctx context.Context, timeout time.Duration) (message.Message, error) {
	return c.connMgr.Recv(ctx, timeout)
}
