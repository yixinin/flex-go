package buffers

import "github.com/yixinin/flex/message"

type Buffer interface {
	Push(msg message.Message)
	Pop() message.Message
	Top() message.Message
	Len() int
}

const (
	BufferQueue = "queue"
)

func NewBufferFunc(name string) func() Buffer {
	switch name {
	case BufferQueue:
		return NewQueue
	}
	panic("no such buffer " + name)
}
