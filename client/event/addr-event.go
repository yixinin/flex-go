package event

import "net"

const (
	EventAdd AddrEventType = 1
	EventDel AddrEventType = 2
)

type AddrEventType uint8

type AddrEvent struct {
	EventType AddrEventType
	Id        string
	Addr      *net.TCPAddr
}
