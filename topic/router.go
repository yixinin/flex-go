package topic

import (
	"context"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type MessageRouter interface {
	Send(ctx context.Context, msg message.Message) error
	OnSubJoin(ctx context.Context, sub *client.Subscriber)
	OnSubLeave(ctx context.Context, id string)
	OnPubJoin(ctx context.Context, pub *client.Publisher)
	OnPubLeave(ctx context.Context, id string)
}

const (
	RouterRoundRobin = "round-robin"
	RouterHash       = "hash"
)

func NewRouter(name string) MessageRouter {
	switch name {
	case RouterRoundRobin:
		return NewRoundRobinRouter()
	case RouterHash:
		return NewHashRouter()
	}
	panic("not such router " + name)
}
