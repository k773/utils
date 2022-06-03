package typedBuffer

import "io"

// OffsetReader returns new reader that replaces ReadAt(offset, n) with just Read(); combine with io.LimitReader to read only n bytes
type OffsetReader struct {
	io.ReaderAt
	Offset int64
}

func (r *OffsetReader) Read(dst []byte) (n int, e error) {
	if r.ReaderAt == nil {
		return 0, io.ErrClosedPipe
	}
	return r.ReadAt(dst, r.Offset)
}
