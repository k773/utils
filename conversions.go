package utils

import (
	"fmt"
	"reflect"
	"unsafe"
)

func Int2strFast(i int) string {
	const bufSize = 20 // Max string length is 20sym (ui64=19sym + minus if neg)
	const bufLastIndex = 19
	var buf = make([]byte, bufSize)
	//var buf []byte

	var (
		n   = bufLastIndex
		neg = i < 0
		v   byte
	)
	if neg {
		i = -i
	}
	for v, i = byte(i%10), i/10; ; n, v, i = n-1, byte(i%10), i/10 {
		buf[n] = '0' + v
		if i == 0 {
			break
		}
	}
	if neg {
		n--
		buf[n] = '-'
	}
	return (*(*string)(unsafe.Pointer(&buf)))[n:]
}

// Uint2strFast converts uint to string
// Uint2strFast(uint(int)) approx 33% faster than Int2strFast (4ns vs 3ns)
func Uint2strFast(i uint) string {
	const bufSize = 20 // Max string length is 20sym (ui64=19sym + minus if neg)
	const bufLastIndex = 20 - 1
	var buf = make([]byte, bufSize)
	//var buf []byte

	var (
		n = bufLastIndex
		v byte
	)
	for v, i = byte(i%10), i/10; ; n, v, i = n-1, byte(i%10), i/10 {
		buf[n] = '0' + v
		if i == 0 {
			break
		}
	}
	return (*(*string)(unsafe.Pointer(&buf)))[n:]
}

func ToString(v any) string {
	var reflV = ReflectDereference(v)

	// Detecting fmt.Stringer, Error()
	var canInterfaceToString = func(v reflect.Value) bool {
		if v.CanInterface() {
			var i = v.Interface()
			if _, ok := i.(fmt.Stringer); ok {
				return ok
			}
			if _, ok := i.(error); ok {
				return ok
			}
		}
		return false
	}

	if !reflV.CanAddr() && !canInterfaceToString(reflV) {
		var newReflV = reflect.New(reflV.Type())
		if canInterfaceToString(newReflV) {
			reflV = newReflV
		}
	}
	return fmt.Sprintf("%v", reflV)
}
