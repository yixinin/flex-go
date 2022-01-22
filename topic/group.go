package topic

import (
	"context"
	"time"

	"github.com/yixinin/flex/buffers"
	"github.com/yixinin/flex/message"
)

type Group struct {
	buffer buffers.Buffer
	TTL    int64
	ch     chan message.Message
}

func newGroup(key string, newBuffer func() buffers.Buffer) *Group {
	g := &Group{
		buffer: newBuffer(),
		TTL:    time.Now().Unix() + 10*60,
	}
	return g
}

func (g *Group) Close() {
	close(g.ch)
}

func (g *Group) Run(ctx context.Context, ch chan message.Message) {
	for msg := range g.ch {
		switch msg := msg.(type) {
		case *message.AckMessage:
			g.buffer.Pop()
			if top := g.buffer.Top(); top != nil {
				ch <- top
			}
		case *message.RawMessage:
			g.buffer.Push(msg)
		}
	}
}
