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

func AddToTheLeft(s string, n int, ch byte) string {
	var mk = make([]byte, n+len(s))
	MemsetRepeat(mk, n, ch)
	copy(mk[n:], s)
	return string(mk)
}
