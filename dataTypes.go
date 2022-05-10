package utils

/*
	Shift-able slice
*/

type ShiftAbleSlice[T comparable] struct {
	Back []T
}

func NewShiftAbleSlice[T comparable]() *ShiftAbleSlice[T] {
	return &ShiftAbleSlice[T]{}
}

func NewShiftAbleSliceFrom[T comparable](a []T) *ShiftAbleSlice[T] {
	return &ShiftAbleSlice[T]{Back: a}
}

func (s *ShiftAbleSlice[T]) ShiftLeft(n int) {
	if len(s.Back) >= n {
		s.Back = s.Back[n:]
	}
}

func (s *ShiftAbleSlice[T]) ShiftRight(n int) {
	if len(s.Back) >= n {
		s.Back = s.Back[:len(s.Back)-n]
	}
}

/*
	If restriction is 0, restriction is not applied
	[..., new-3, new-2, new-1, new]
*/
type RestrictedShiftAbleSliceRight[T comparable] struct {
	*ShiftAbleSlice[T]
	RestrictedSize int
}

func NewRestrictedShiftAbleSliceRight[T comparable](restriction int, preallocate bool) *RestrictedShiftAbleSliceRight[T] {
	if preallocate {
		var a = make([]T, 0, restriction)
		return &RestrictedShiftAbleSliceRight[T]{NewShiftAbleSliceFrom(a), restriction}
	} else {
		return &RestrictedShiftAbleSliceRight[T]{NewShiftAbleSlice[T](), restriction}
	}
}

func (r *RestrictedShiftAbleSliceRight[T]) Put(a T) {
	var d = len(r.Back) - r.RestrictedSize
	if d >= 0 {
		r.ShiftLeft(d + 1)
	}
	r.Back = append(r.Back, a)
}

func (r *RestrictedShiftAbleSliceRight[T]) SetLast(a T) {
	if len(r.Back) == 0 {
		r.Back = append(r.Back, a)
	} else {
		r.Back[len(r.Back)-1] = a
	}
}

/*
	If restriction is 0, restriction is not applied
	[new, new-1, new-2, new-3, ...]
*/
type RestrictedShiftAbleSliceLeft[T comparable] struct {
	*ShiftAbleSlice[T]
	RestrictedSize int
}

func NewRestrictedShiftAbleSliceLeft[T comparable](restriction int, preallocate bool) *RestrictedShiftAbleSliceLeft[T] {
	if preallocate {
		var a = make([]T, 0, restriction)
		return &RestrictedShiftAbleSliceLeft[T]{NewShiftAbleSliceFrom(a), restriction}
	} else {
		return &RestrictedShiftAbleSliceLeft[T]{NewShiftAbleSlice[T](), restriction}
	}
}

func (r *RestrictedShiftAbleSliceLeft[T]) Put(a T) {
	var d = len(r.Back) - r.RestrictedSize
	if d >= 0 {
		r.ShiftRight(d + 1)
	}
	r.Back = append([]T{a}, r.Back...)
}

func (r *RestrictedShiftAbleSliceLeft[T]) SetLast(a T) {
	if len(r.Back) == 0 {
		r.Back = append(r.Back, a)
	} else {
		r.Back[0] = a
	}
}
