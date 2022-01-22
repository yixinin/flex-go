package topic

import (
	"context"
	"time"

	"github.com/yixinin/flex/buffers"
	"github.com/yixinin/flex/logger"
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
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-g.ch:
			switch msg := msg.(type) {
			case *message.AckMessage:
				pop := g.buffer.Pop()
				if top := g.buffer.Top(); top != nil {
					ch <- top
				}
				logger.Debugf(ctx, "message:%+v acked", pop)
			case *message.RawMessage:
				if g.buffer.Len() == 0 {
					ch <- msg
				}
				g.buffer.Push(msg)
			}
		}
	}
}
