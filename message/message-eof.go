package message

type Eof struct {
}

func (m Eof) Id() string {
	return ""
}
func (m Eof) Group() string {
	return ""
}
func (m Eof) RawData() []byte {
	return nil
}

func (m Eof) PeerId() string {
	return ""
}
