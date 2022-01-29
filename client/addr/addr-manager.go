package addr

import (
	"context"

	"github.com/yixinin/flex/client/config"
	"github.com/yixinin/flex/client/event"
)

type AddrManager interface {
	Event() chan event.AddrEvent
	Run(ctx context.Context) error
}

type NewAddrManager func(conf config.Config) AddrManager
