package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Big struct {
	A, B uint64
}

func BigZero() *Big {
	var zero Big
	return &zero
}

func (b *Big) Clone() *Big {
	return &Big{A: b.A, B: b.B}
}

func (b *Big) Clear() {
	b.A, b.B = 0, 0
}

func (b *Big) Or(b2 *Big) *Big {
	b.A |= b2.A
	b.B |= b2.B
	return b
}

func (b *Big) SetBit(n int, a uint) *Big {
	if n <= 63 {
		b.A |= (uint64(a) & 1) << n
	} else {
		b.B |= (uint64(a) & 1) << (n - 64)
	}
	return b
}

func (b *Big) And(b2 *Big) *Big {
	b.A &= b2.A
	b.B &= b2.B
	return b
}

func (b *Big) ShiftRight(bits int) *Big {
	b.A >>= bits
	for i := 0; i < bits; i++ {
		var x uint64 = 1 << i
		if b.B&x != 0 {
			b.A |= 1 << (63 - i)
		}
	}
	b.B >>= bits

	return b
}

func (b *Big) ShiftLeft(bits int) *Big {
	b.B <<= bits
	for i := 0; i <= bits; i++ {
		var x uint64 = 1 << (64 - i)
		if b.A&x != 0 {
			b.B |= 1 << (bits - i)
		}
	}
	b.A <<= bits

	return b
}

func (b *Big) Equals(b2 *Big) bool {
	return b.A == b2.A && b.B == b2.B
}

func (b *Big) IsZero() bool {
	return b.A == 0 && b.B == 0
}

func (b *Big) SetUInt64(a0, a1 uint64) *Big {
	b.A = a0
	b.B = a1
	return b
}

func (b *Big) String() string {
	var res = make([]byte, 128)
	binary.BigEndian.PutUint64(res, b.A)
	binary.BigEndian.PutUint64(res[64:], b.B)
	return hex.EncodeToString(res)
}

func (b *Big) FromString(s string) *Big {
	a, e := hex.DecodeString(s)
	if e != nil || len(a) != 128 {
		b.Clear()
	} else {
		b.A = binary.BigEndian.Uint64(a)
		b.B = binary.BigEndian.Uint64(a[64:])
	}
	return b
}

func (b *Big) FormatBits() string {
	return fmt.Sprintf("%064b %064b", b.B, b.A)
}

// Fail

//
//import (
//	"encoding/binary"
//	"encoding/hex"
//	"reflect"
//	"unsafe"
//)
//
//// Code for comparison is taken from here: https://github.com/grailbio/base/blob/master/simd/and_amd64.go
//
//// BytesPerWord is the number of bytes in a machine word.
//// We don't use unsafe.Sizeof(uintptr(1)) since there are advantages to having
//// this as an untyped constant, and there's essentially no drawback since this
//// is an _amd64-specific file.
//const BytesPerWord = 8
//
//// Log2BytesPerWord is log2(BytesPerWord).  This is relevant for manual
//// bit-shifting when we know that's a safe way to divide and the compiler does
//// not (e.g. dividend is of signed int type).
//const Log2BytesPerWord = uint(3)
//
//// AndUnsafeInplace sets main[pos] := main[pos] & arg[pos] for every position
//// in main[].
////
//// WARNING: This is a function designed to be used in inner loops, which makes
//// assumptions about length and capacity which aren't checked at runtime.  Use
//// the safe version of this function when that's a problem.
//// Assumptions #2-3 are always satisfied when the last
//// potentially-size-increasing operation on arg[] is {Re}makeUnsafe(),
//// ResizeUnsafe(), or XcapUnsafe(), and the same is true for main[].
////
//// 1. len(arg) and len(main) must be equal.
////
//// 2. Capacities are at least RoundUpPow2(len(main) + 1, bytesPerVec).
////
//// 3. The caller does not care if a few bytes past the end of main[] are
//// changed.
//func AndUnsafeInplace(main, arg []byte) {
//	mainLen := len(main)
//	argHeader := (*reflect.SliceHeader)(unsafe.Pointer(&arg))
//	mainHeader := (*reflect.SliceHeader)(unsafe.Pointer(&main))
//	argWordsIter := unsafe.Pointer(argHeader.Data)
//	mainWordsIter := unsafe.Pointer(mainHeader.Data)
//	if mainLen > 2*BytesPerWord {
//		nWordMinus2 := (mainLen - BytesPerWord - 1) >> Log2BytesPerWord
//		for widx := 0; widx < nWordMinus2; widx++ {
//			mainWord := *((*uintptr)(mainWordsIter))
//			argWord := *((*uintptr)(argWordsIter))
//			*((*uintptr)(mainWordsIter)) = mainWord & argWord
//			mainWordsIter = unsafe.Pointer(uintptr(mainWordsIter) + BytesPerWord)
//			argWordsIter = unsafe.Pointer(uintptr(argWordsIter) + BytesPerWord)
//		}
//	} else if mainLen <= BytesPerWord {
//		mainWord := *((*uintptr)(mainWordsIter))
//		argWord := *((*uintptr)(argWordsIter))
//		*((*uintptr)(mainWordsIter)) = mainWord & argWord
//		return
//	}
//	// The last two read-and-writes to main[] usually overlap.  To avoid a
//	// store-to-load forwarding slowdown, we read both words before writing
//	// either.
//	// shuffleLookupOddInplaceSSSE3Asm() uses the same strategy.
//	mainWord1 := *((*uintptr)(mainWordsIter))
//	argWord1 := *((*uintptr)(argWordsIter))
//	finalOffset := uintptr(mainLen - BytesPerWord)
//	mainFinalWordPtr := unsafe.Pointer(mainHeader.Data + finalOffset)
//	argFinalWordPtr := unsafe.Pointer(argHeader.Data + finalOffset)
//	mainWord2 := *((*uintptr)(mainFinalWordPtr))
//	argWord2 := *((*uintptr)(argFinalWordPtr))
//	*((*uintptr)(mainWordsIter)) = mainWord1 & argWord1
//	*((*uintptr)(mainFinalWordPtr)) = mainWord2 & argWord2
//}
//
//func OrUnsafeInplace(main, arg []byte) {
//	mainLen := len(main)
//	argHeader := (*reflect.SliceHeader)(unsafe.Pointer(&arg))
//	mainHeader := (*reflect.SliceHeader)(unsafe.Pointer(&main))
//	argWordsIter := unsafe.Pointer(argHeader.Data)
//	mainWordsIter := unsafe.Pointer(mainHeader.Data)
//	if mainLen > 2*BytesPerWord {
//		nWordMinus2 := (mainLen - BytesPerWord - 1) >> Log2BytesPerWord
//		for widx := 0; widx < nWordMinus2; widx++ {
//			mainWord := *((*uintptr)(mainWordsIter))
//			argWord := *((*uintptr)(argWordsIter))
//			*((*uintptr)(mainWordsIter)) = mainWord | argWord
//			mainWordsIter = unsafe.Pointer(uintptr(mainWordsIter) + BytesPerWord)
//			argWordsIter = unsafe.Pointer(uintptr(argWordsIter) + BytesPerWord)
//		}
//	} else if mainLen <= BytesPerWord {
//		mainWord := *((*uintptr)(mainWordsIter))
//		argWord := *((*uintptr)(argWordsIter))
//		*((*uintptr)(mainWordsIter)) = mainWord | argWord
//		return
//	}
//	// The last two read-and-writes to main[] usually overlap.  To avoid a
//	// store-to-load forwarding slowdown, we read both words before writing
//	// either.
//	// shuffleLookupOddInplaceSSSE3Asm() uses the same strategy.
//	mainWord1 := *((*uintptr)(mainWordsIter))
//	argWord1 := *((*uintptr)(argWordsIter))
//	finalOffset := uintptr(mainLen - BytesPerWord)
//	mainFinalWordPtr := unsafe.Pointer(mainHeader.Data + finalOffset)
//	argFinalWordPtr := unsafe.Pointer(argHeader.Data + finalOffset)
//	mainWord2 := *((*uintptr)(mainFinalWordPtr))
//	argWord2 := *((*uintptr)(argFinalWordPtr))
//	*((*uintptr)(mainWordsIter)) = mainWord1 | argWord1
//	*((*uintptr)(mainFinalWordPtr)) = mainWord2 | argWord2
//}
//
//type Big [16]byte
//
//func BigZero() *Big {
//	var zero Big
//	return &zero
//}
//
//func (b *Big) Clone() *Big {
//	var b2 Big
//	copy(b2[:], b[:])
//	return &b2
//}
//
//func (b *Big) Clear() {
//	var empty Big
//	*b = empty
//}
//
//func (b *Big) Or(b2 *Big) *Big {
//	OrUnsafeInplace(b[:], b2[:])
//	return b
//}
//
//func (b *Big) SetBit(n int, a bool) *Big {
//	d := 0
//	if a {
//		d = 1
//	}
//	e := byte(n % 8)
//
//	b[n/8] = b[n/8]&^(1<<e) | byte(d<<e)
//	return b
//}
//
//func (b *Big) And(b2 *Big) *Big {
//	AndUnsafeInplace(b[:], b2[:])
//	//for i := range b {
//	//	b[i] &= b2[i]
//	//}
//	return b
//}
//
//func (b *Big) ShiftRight(bits int) *Big {
//	ShiftLeft(b[:], -bits)
//
//	return b
//}
//
//func (b *Big) ShiftLeft(bits int) *Big {
//	ShiftLeft(b[:], bits)
//
//	return b
//}
//
//func (b *Big) Equals(b2 *Big) bool {
//	return string(b[:]) == string(b2[:])
//}
//
//func (b *Big) IsZero() bool {
//	var empty Big
//	return b.Equals(&empty)
//}
//
//func (b *Big) SetUInt64(v uint64) *Big {
//	b.Clear()
//	binary.BigEndian.PutUint64(b[:], v)
//	return b
//}
//
//func (b *Big) UInt64() uint64 {
//	return binary.BigEndian.Uint64(b[:])
//}
//
//func (b *Big) String() string {
//	return hex.EncodeToString(b[:])
//}
//
//func (b *Big) FromString(s string) *Big {
//	a, e := hex.DecodeString(s)
//	if e != nil || len(a) != len(b) {
//		b.Clear()
//	} else {
//		copy((*b)[:], a)
//	}
//	return b
//}
