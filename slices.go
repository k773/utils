package utils

/*
	Slice length reported
*/

type SliceLengthReporterInterface interface {
	Len() int
}

type SliceLengthReporter[T any] struct {
	Slice *[]T
}

func (s *SliceLengthReporter[T]) Len() int {
	return len(*s.Slice)
}
