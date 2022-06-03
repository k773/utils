package readerWriterAt

type SimpleReaderWriterAt struct {
	Buf []byte
}

func NewSimpleReaderWriterAt(buf []byte) *SimpleReaderWriterAt {
	return &SimpleReaderWriterAt{Buf: buf}
}

func (s *SimpleReaderWriterAt) WriteAt(p []byte, off int64) (n int, err error) {
	var max = int(off) + len(p)
	s.Grow(max)
	return copy(s.Buf[off:], p), nil
}

func (s *SimpleReaderWriterAt) ReadAt(p []byte, off int64) (n int, err error) {
	return copy(p, s.Buf[off:]), nil
}

func (s *SimpleReaderWriterAt) Grow(n int) {
	if n == 0 {
		s.Reset()
		return
	}

	if len(s.Buf) >= n {
		return
	}
	// Trying to grow by resize
	if cap(s.Buf) >= n {
		s.Buf = s.Buf[:n]
		return
	}
	// Growing by re-making
	buf := make([]byte, n)
	copy(buf, s.Buf)
	s.Buf = buf
}

func (s *SimpleReaderWriterAt) Reset() {
	s.Buf = s.Buf[:0]
}
