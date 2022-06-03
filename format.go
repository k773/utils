package utils

import (
	"strconv"
	"strings"
)

/*
	Bytes
*/

type BytesSize int

func (b BytesSize) String() string {
	return FormatBytesSize(int(b), 2)
}

func FormatBytesSize(b, precision int) string {
	switch {
	case b < 1<<10:
		return strconv.Itoa(b) + " B"
	case b < 1<<20:
		return strconv.FormatFloat(float64(b)/float64(1<<10), 'f', precision, 64) + " KB"
	case b < 1<<30:
		return strconv.FormatFloat(float64(b)/float64(1<<20), 'f', precision, 64) + " MB"
	case b < 1<<40:
		return strconv.FormatFloat(float64(b)/float64(1<<30), 'f', precision, 64) + " GB"
	case b < 1<<50:
		return strconv.FormatFloat(float64(b)/float64(1<<40), 'f', precision, 64) + " TB"
	case b < 1<<60:
		return strconv.FormatFloat(float64(b)/float64(1<<50), 'f', precision, 64) + " PB"
	default:
		return strconv.FormatFloat(float64(b)/float64(1<<60), 'f', precision, 64) + " EB"
	}
}

/*
	Bits
*/

type BitsSize int

func (b BitsSize) String() string {
	return FormatBitsSize(int(b), 2)
}

func FormatBitsSize(b, precision int) string {
	a := FormatBytesSize(b, precision)
	return a[:len(a)-1] + strings.ToLower(string(a[len(a)-1]))
}
