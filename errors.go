package utils

type TemporaryError struct {
	msg string
}

func (e *TemporaryError) Error() string {
	return e.msg
}

func (e *TemporaryError) Temporary() bool {
	return true
}

func NewTemporaryError(s string) *TemporaryError {
	return &TemporaryError{msg: s}
}
