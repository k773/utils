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
