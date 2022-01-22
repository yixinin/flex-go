package topic

import (
	"context"

	"github.com/yixinin/flex/client"
	"github.com/yixinin/flex/message"
)

type ReturnSender struct {
	publishers map[string]*client.Publisher
}

func (m *ReturnSender) Send(ctx context.Context, msg message.Message) (err error) {
	if pub, ok := m.publishers[msg.ClientId()]; ok {
		return pub.Send(ctx, msg)
	}
	return nil
}
