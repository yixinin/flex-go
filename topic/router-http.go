package topic

import (
	"context"
	"net/http"

	"github.com/yixinin/flex/message"
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
