package client

import (
	"context"
	"net"
)

const (
	EventAdd = 1
	EventDel = 2
)

type Event struct {
	EventType int
	Id        string
	Addr      *net.TCPAddr
}

func (c *Client) recvEvent(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-c.event:
			switch ev.EventType {
			case EventAdd:
				c.connMgr.OnAddrJoin(ctx, ev.Id, ev.Addr)
			case EventDel:
				c.connMgr.OnAddrLeave(ctx, ev.Id)
			}
		}
	}
}
