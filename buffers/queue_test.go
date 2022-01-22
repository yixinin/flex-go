package buffers

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/yixinin/flex/message"
)

func TestQueue(t *testing.T) {
	var queue = NewQueue()
	for i := 0; i < 10; i++ {
		queue.Push(&message.AckMessage{
			Key: strconv.Itoa(i),
		})
		if i%2 == 0 {
			fmt.Println(queue.Pop())
		}

	}

	msg := queue.Pop()
	for msg != nil {
		msg = queue.Pop()
		fmt.Println(msg)
	}
}
