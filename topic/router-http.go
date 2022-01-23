package topic

import (
	"context"
	"net/http"

	"github.com/yixinin/flex/message"
	"github.com/yixinin/flex/pubsub"
)

type HttpSender struct {
	url    string
	method string
	client http.Client
}

func (m *HttpSender) Send(ctx context.Context, msg message.Message) (err error) {
	req := &http.Request{
		Method: m.method,
	}
	_, err = m.client.Do(req)
	return err
}

func (m *HttpSender) OnSubJoin(ctx context.Context, sub *pubsub.Subscriber) {}
func (m *HttpSender) OnSubLeave(ctx context.Context, id string)             {}
func (m *HttpSender) OnPubJoin(ctx context.Context, pub *pubsub.Publisher)  {}
func (m *HttpSender) OnPubLeave(ctx context.Context, id string)             {}
