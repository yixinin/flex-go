package buffers

import (
	"github.com/yixinin/flex/message"
)

type Node struct {
	next *Node
	prev *Node
	val  message.Message
}

type Queue struct {
	head *Node
	tail *Node
	size int
}

func NewQueue() Buffer {
	return &Queue{}
}

func (q *Queue) Push(msg message.Message) {
	q.size++
	node := &Node{
		val: msg,
	}
	if q.head == nil {
		q.head = node
		q.tail = node
		return
	}
	head := q.head
	q.head = node
	head.prev = node
	node.next = head
}
func (q *Queue) Pop() message.Message {
	tail := q.tail
	if tail != nil {
		q.tail = tail.prev
		if q.tail == nil {
			q.head = nil
		}
		q.size--
		return tail.val
	}

	return nil
}
func (q *Queue) Top() message.Message {
	tail := q.tail
	if tail != nil {
		return tail.val
	}

	return nil
}
func (q *Queue) Len() int {
	return q.size
}
