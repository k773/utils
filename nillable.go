package utils

import "errors"

type Nillable[T any] struct {
	ValuePtr *T
}

func (n *Nillable[T]) IsNil() bool {
	return n.ValuePtr == nil
}

// With calls do with Nillable.Value if the value is not nil.
// If the Nillable.Value is nil, an error is returned.
func (n *Nillable[T]) With(do func(T) error) error {
	return n.WithOrError(do, errors.New("nil value"))
}

// WithOrError calls do with Nillable.Value if the value is not nil.
// If the Nillable.Value is nil, the passed ifNil error is returned.
func (n *Nillable[T]) WithOrError(do func(T) error, ifNil error) error {
	var v = n.ValuePtr
	if v == nil {
		return ifNil
	}
	return do(*v)
}
