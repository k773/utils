package utils

import "encoding/binary"

// TODO: make it work
//func ShiftLeft(data []byte, bits int) {
//	n := len(data)
//	if bits < 0 {
//		bits = -bits
//		for i := n - 1; i > 0; i-- {
//			data[i] = data[i]>>bits | data[i-1]<<(8-bits)
//		}
//		data[0] >>= bits
//	} else {
//		for i := 0; i < n-1; i++ {
//			data[i] = data[i]<<bits | data[i+1]>>(8-bits)
//		}
//		data[n-1] <<= bits
//	}
//}

func Uint64ToLE(a uint64) []byte {
	var buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, a)
	return buf
}
func Uint64ToBE(a uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, a)
	return buf
}
