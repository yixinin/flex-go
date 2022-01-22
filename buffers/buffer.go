package buffers

import "github.com/yixinin/flex/message"

type Buffer interface {
	Push(msg message.Message)
	Pop() message.Message
	Top() message.Message
	Closed() bool
}
