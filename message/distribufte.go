package message

type DistributeType byte

const (
	DisNone     DistributeType = 0
	Return      DistributeType = 1
	HttpRequest DistributeType = 2
	Hash        DistributeType = 3
	RoundRobin  DistributeType = 4
)
