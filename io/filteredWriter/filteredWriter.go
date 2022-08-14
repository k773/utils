package filteredWriter

import (
	"bytes"
	"io"
)

type FilteredWriter struct {
	Buf *bytes.Buffer
	Dst io.Writer

	// Filter should return true if the data should be written to the Dst. In other case, the data will be discarded.
	Filter func(data []byte) bool
	// FilterDelimiter defines the byte up to which (inclusive) all data will be passed to the Filter.
	FilterDelimiter byte

	// OverflowProtectionMaxSize defines the max Buf size before dumping all data to the Dst. Set to 0 to remove the limit.
	OverflowProtectionMaxSize int
}

func NewFilteredWriter(dst io.Writer, filter func(data []byte) bool, filterDelimiter byte, overflowProtectionMaxSize int) *FilteredWriter {
	return &FilteredWriter{Buf: bytes.NewBuffer(nil), Dst: dst, Filter: filter, FilterDelimiter: filterDelimiter, OverflowProtectionMaxSize: overflowProtectionMaxSize}
}

func (f *FilteredWriter) Write(p []byte) (n int, err error) {
	f.Buf.Write(p)

	var data []byte
	for err == nil {
		data, err = f.Buf.ReadBytes(f.FilterDelimiter)
		if err == nil {
			if f.Filter(data) {
				_, err = f.Dst.Write(data)
			}
		}
	}
	// Returning the data to the buffer
	if len(data) != 0 {
		f.Buf.Write(data)
	}
	// Overflow protection
	if f.OverflowProtectionMaxSize != 0 && f.Buf.Len() >= f.OverflowProtectionMaxSize {
		_, err = io.Copy(f.Dst, f.Buf)
	}
	n = len(p)
	return
}
