package utils

import (
	"io"
	"unsafe"
)

func RandStringMust(reader io.Reader, length int) (res string) {
	res, e := RandString(reader, length)
	if e != nil {
		panic(e)
	}
	return res
}

// RandString generates random string of characters [0-9a-Z]
func RandString(reader io.Reader, length int) (res string, e error) {
	// 10 + 26 + 26 = 62
	var getN = func(r uint8) (v uint8) {
		if r >= 10 {
			r += 7
			if r >= 43 {
				r += 6
			}
		}
		return r + '0'
	}

	var buf = make([]byte, length)
	if _, e = reader.Read(buf); e == nil {
		for i, v := range buf {
			buf[i] = getN(v % 62)
		}
	}
	return string(buf), e
}

func RandomChoice[T any](reader io.Reader, src []T) (T, error) {
	var v T
	var r, e = RandValue[uint](reader)
	if e == nil {
		v = src[r%uint(len(src))]
	}
	return v, e
}

func RandomChoiceMulti[T any](reader io.Reader, src []T, resLength int) ([]T, error) {
	var v = make([]T, resLength)
	var r, e = RandArray[uint](reader, make([]uint, resLength))
	if e == nil {
		for i, rv := range r {
			v[i] = src[rv%uint(len(src))]
		}
	}
	return v, e
}

func RandArray[T any](reader io.Reader, arr []T) ([]T, error) {
	var v T
	var tmp2 = unsafe.Slice(&arr[0], len(arr)*int(unsafe.Sizeof(v)))
	var tmp = *(*[]byte)(unsafe.Pointer(&tmp2))

	_, e := reader.Read(tmp)

	return arr, e
}

func RandValue[T any](reader io.Reader) (T, error) {
	var v T
	var tmp2 = unsafe.Slice(&v, int(unsafe.Sizeof(v)))
	var tmp = *(*[]byte)(unsafe.Pointer(&tmp2))

	_, e := reader.Read(tmp)

	return v, e
}
