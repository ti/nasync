package nasync

const (
	ERROR_UNKNOWN = iota
	ERROR_TIMEOUT
	ERROR_QUEUE_FULL
	ERROR_QUEUE_CLOSED
)

type Error struct {
	Code int

	Msg string
}

func (e *Error) Error() string {

	return e.Msg
}
