package utils

func PadLeft(s string, n int, ch byte) string {
	l := len(s)
	if l < n {
		var d = n - l
		var mk = make([]byte, n)
		copy(mk[d:], s)
		MemsetRepeat(mk, d, ch)
		return string(mk)
	}
	return s
}
