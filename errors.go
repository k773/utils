package utils

import "errors"

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

type ComparableTextError struct {
	Text string
}

func (c *ComparableTextError) Is(err error) bool {
	var cast *ComparableTextError
	return errors.As(err, &cast) && cast.Text == c.Text
}

func (c *ComparableTextError) Error() string {
	return c.Text
}
