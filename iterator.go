package utils

func IterateF[srcT any, dstT any](src []srcT, iter func(dstT), convert func(srcT) dstT) {
	for _, v := range src {
		iter(convert(v))
	}
}

// IterateFCond iterates over given src until either eof or convert() returns false
func IterateFCond[srcT any, dstT any](src []srcT, iter func(dstT) bool, convert func(srcT) dstT) {
	for _, v := range src {
		if !iter(convert(v)) {
			break
		}
	}
}

type IteratorF[srcT any, dstT any] struct {
	L       []srcT
	Convert func(srcT) dstT
}

func (rcv *IteratorF[srcT, dstT]) IterateF(iter func(dstT)) {
	for _, v := range rcv.L {
		iter(rcv.Convert(v))
	}
}

// IterateFCond iterates over given src until either eof or convert() returns false
func (rcv *IteratorF[srcT, dstT]) IterateFCond(iter func(dstT) bool) {
	for _, v := range rcv.L {
		if !iter(rcv.Convert(v)) {
			break
		}
	}
}

/*
	Simple 'foreach' iterator
*/

func ForEachSlice[T any](slice []T, f func(v T) (wantMore bool), reverse bool) {
	if !reverse {
		for _, v := range slice {
			if !f(v) {
				break
			}
		}
	} else {
		for i := len(slice) - 1; i >= 0; i-- {
			if !f(slice[i]) {
				break
			}
		}
	}
}
