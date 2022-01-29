package conn

import (
	"context"
	"net"
	"time"

	"github.com/yixinin/flex/client/config"
	"github.com/yixinin/flex/message"
)

type ConnManager interface {
	OnAddrJoin(ctx context.Context, id string, addr *net.TCPAddr)
	OnAddrLeave(ctx context.Context, id string)
	Recv(ctx context.Context, timeout time.Duration) (message.Message, error)
	Send(ctx context.Context, key, groupKey string, payload []byte) error
	SendAsync(ctx context.Context, key, groupKey string, payload []byte)
	Run(ctx context.Context) error
}

type NewConnManager func(conf config.Config) ConnManager
