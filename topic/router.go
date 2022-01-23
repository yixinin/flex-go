package topic

import (
	"context"

	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
)

type MessageRouter interface {
	Send(ctx context.Context, msg message.Message) error
	OnSubJoin(ctx context.Context, sub *pubsub.Subscriber)
	OnSubLeave(ctx context.Context, id string)
	OnPubJoin(ctx context.Context, pub *pubsub.Publisher)
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
