package pubsub

import (
	"net"

	"github.com/yixinin/flex/message"
)

type Subscriber struct {
	*Client
}

func NewSubscriber(conn net.Conn, msg *message.ConnMessage, cancel func()) *Subscriber {
	return &Subscriber{
		Client: newClient(conn, msg, cancel),
	}
}
