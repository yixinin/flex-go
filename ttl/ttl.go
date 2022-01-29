package ttl

import "time"

func NextTTL() int64 {
	return time.Now().Add(time.Second).UnixMilli()
}

func NewTicker() *time.Ticker {
	return time.NewTicker(500 * time.Millisecond)
}
