package client

import (
	"net"

	"github.com/yixinin/flex/message"
)

type Publisher struct {
	*Client
}

func NewPublisher(conn net.Conn, msg *message.ConnMessage, cancel func()) *Publisher {
	return &Publisher{
		Client: newClient(conn, msg, cancel),
	}
}
