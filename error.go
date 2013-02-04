package netchan

type ClosedError struct {}

func (err *ClosedError) Error() string {
	return "Channel is closed"
}

var (
	Closed = new(ClosedError)
)
