package domain

type Direction string

const (
	SendDirection    Direction = "send"
	ReceiveDirection Direction = "receive"
)

func (d Direction) Valid() bool {
	switch d {
	case SendDirection, ReceiveDirection:
		return true
	default:
		return false
	}
}

func (d Direction) String() string {
	return string(d)
}
