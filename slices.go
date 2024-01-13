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

func SliceFilter[T any](in []T, keep func(T) bool) []T {
	var res []T
	for _, v := range in {
		if keep(v) {
			res = append(res, v)
		}
	}
	return res
}
