package client

import (
	"context"
	"net"
	"time"

	"github.com/yixinin/flex/client/event"
	"github.com/yixinin/flex/message"
)

type AddrManager interface {
	Event() chan event.AddrEvent
	Run(ctx context.Context) error
}

type ConnManager interface {
	OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr)
	OnAddrLeave(ctx context.Context, id string)
	Recv(ctx context.Context, timeout time.Duration) (message.Message, error)
	Send(ctx context.Context, key, groupKey string, payload []byte) error
	SendAsync(ctx context.Context, key, groupKey string, payload []byte)
	Ack(ctx context.Context, clientId, key, groupKey string)
	Run(ctx context.Context) error
}
