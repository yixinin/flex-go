package message

type CloseMessage struct {
}

func NewCloseMessage() Message {
	return &CloseMessage{}
}
