package utils

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
