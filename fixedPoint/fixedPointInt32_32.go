/*

 */

package fixedPoint

import (
	"github.com/k773/utils"
	"strconv"
)

type Int32_32 int64

func NewInt3232(decimal int32, fractional uint32, fractionalPow uint8) Int32_32 {
	const (
		powShift = 32 - 5
	)

	return Int32_32(decimal)<<32 | Int32_32(1<<5-1-fractionalPow)<<powShift | Int32_32(fractional)
}

func ParseFixedPointToInt3232(src string) Int32_32 {
	a, b, c, _ := ParseFixedPointRaw(src, -1)
	return NewInt3232(int32(a), uint32(b), uint8(c))
}

func (i Int32_32) Parts() (decimal int32, fractional uint32, fractionalPow uint8) {
	const (
		powShift       = 32 - 5
		powBits        = 1<<5 - 1
		powMask        = powBits << powShift
		fractionalMask = ^powMask
	)

	return int32(i >> 32), uint32(i & fractionalMask), uint8(powBits - (i & powMask >> powShift))
}

func (i Int32_32) String() string {
	decimal, fractional, fractionalPow := i.Parts()
	return strconv.Itoa(int(decimal)) + "." + utils.AddToTheLeft(strconv.Itoa(int(fractional)), int(fractionalPow), '0')
}

//func (i Int32_32) Add(i2 Int32_32) Int32_32 {
//	const (
//		decimalShift   = 32
//		powBitsLen     = 5
//		powShift       = decimalShift - powBitsLen
//		decimalMask    = ^(1<<decimalShift - 1)
//		powBits        = 1<<powBitsLen - 1
//		powMask        = powBits << powShift
//		fractionalMask = 1<<(32-5) - 1
//		last32BitsMask = 1<<32 - 1
//	)
//	fmt.Printf("%064b (src) (%v)\n", i, i.String())
//	fmt.Printf("%064b (src) (%v)\n", i2, i2.String())
//
//	var a1, b1, c1 = i.Parts()
//	var a2, b2, c2 = i2.Parts()
//
//	var e1 = utils.Log2(uint64(b1))
//	var e2 = utils.Log2(uint64(b2))
//	var e3 = utils.If(e1 > e2, e1, e2)
//	var f = utils.If(e1 > e2, i&powMask, i2&powMask)
//
//	_ = c1
//	_ = c2
//
//	var d1 = uint64(a1)<<32 | uint64(b1)<<(31-e3-int(c1))
//	var d2 = uint64(a2)<<32 | uint64(b2)<<(31-e3-int(c2))
//
//	fmt.Printf("%064b (%v)\n", d1, i.String())
//	fmt.Printf("%064b (%v)\n", d2, i2.String())
//
//	var d3 = int64(d1) + int64(d2)
//	fmt.Printf("%064b\n", d3)
//	d3 = (d3 & decimalMask) | ((d3 & last32BitsMask) >> (31 - e3)) | int64(f)
//	fmt.Printf("%064b\n", d3)
//
//	return Int32_32(d3)
//
//	//fmt.Println(i.String())
//	//
//	////fmt.Printf("%064b\n", ^uint64(uint32(1<<32-1)))
//	////fmt.Printf("%064b\n", i&decimalMask>>32)
//	////fmt.Printf("%064b %v\n", int32(i&decimalMask>>32+i2&decimalMask>>32), int32(i&decimalMask>>32+i2&decimalMask>>32))
//	//
//	////a, b, c := i.Parts()
//	////a2, b2, c2 := i2.Parts()
//	////
//	////as := a + a2
//	////fixed.Int52_12()
//	////
//	//
//	////return NewInt3232(as, b+b2, (c+c2)/2)
//	//
//	//// Attempting to get a real number
//	//var aLog2 = utils.Log2(uint64(i & fractionalMask))
//	//var a = (i & decimalMask >> 5) | (i&fractionalMask)<<(31-5-aLog2)
//	//
//	//var bLog2 = utils.Log2(uint64(i2 & fractionalMask))
//	//var b = (i2 & decimalMask >> 5) | (i2&fractionalMask)<<(31-5-bLog2)
//	//_ = b
//	//
//	////fmt.Printf("decimal: %064b\nfractional: %064b\nlog2fractional: %v\n", i&fractionalMask, i&decimalMask>>5, aLog2)
//	////fmt.Printf("%064b\n", i&fractionalMask)
//	////fmt.Printf("%064b\n", i&decimalMask>>5)
//	////fmt.Printf("%064b\n", a)
//	//fmt.Printf("%064b\n", a)
//	//fmt.Printf("%064b\n", b)
//	//fmt.Printf("%064b\n", i)
//	//
//	//var c = a + b
//	//var d = aLog2
//	//_ = d
//	//
//	//fmt.Printf("%064b\n", c)
//	////fmt.Printf("%064b\n", c&(fractionalMask)>>(31-5-d))
//	////fmt.Printf("%064b\n", fractionalMask)
//	////fmt.Printf("%064b\n", c&(decimalMask>>5)<<5)
//	////fmt.Printf("%064b\n", c)
//	//c = (c & (decimalMask >> 5) << 5) | (i & powMask) | (c & (fractionalMask) >> (31 - 5 - d))
//	//fmt.Printf("%064b\n", c)
//	//
//	//return c
//}
