package topic

import (
	"context"

	"github.com/yixinin/flex/message"
)

type MessageRouter interface {
	Send(ctx context.Context, msg message.Message) error
}
