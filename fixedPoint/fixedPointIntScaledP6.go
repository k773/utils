/*
	FP implements fixed point math with precision up to 6 decimal places after the dot
*/

package fixedPoint

import (
	"github.com/k773/utils"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"math"
	"unsafe"
)

type FP int64
type IntScaledP6 = FP

// Constants
// If you want to change those, duplicate the type.
const intScaledP6N = 6
const intScaledP6Scale = 10 * 10 * 10 * 10 * 10 * 10
const NaN = -(1<<63 - 1) // bits set are identical to uint64(1<<64-1)

// These are here to remove the unnecessary calculations

// Quick definitions:
const (
	Zero          FP = 0
	ZeroPointOne  FP = 100000
	ZeroPointFive FP = 500000
	One           FP = 1000000
)

func Parse(src string) FP {
	a, b, c, n := ParseFixedPointRaw(src, intScaledP6N)
	return New(a, b, c, n)
}

// ParseBytes parses a fixed point number from a byte slice.
// The byte slice is parsed without allocations that occur when converting to a string.
func ParseBytes(src []byte) FP {
	a, b, c, n := ParseFixedPointRaw(*(*string)(unsafe.Pointer(&src)), intScaledP6N)
	return New(a, b, c, n)
}

// New creates a new FP object.
//
// Arguments:
// - decimal - the number before the dot;
// - fractional - the number after the dot;
// - leadingZeroes - the number of leading zeroes in the fractional part.
// - negative - if the number is negative.
//
// Examples:
// - New(2, 1, 0, false).String() = "2.1"
// - New(2, 1, 1, false).String() = "2.01"
// - New(2, 1, 2, false).String() = "2.001"
func New(decimal, fractional, leadingZeroes uint, negative bool) FP {
	var a = FP(int64(decimal)*intScaledP6Scale + (int64(fractional) * int64(utils.Pow10(int(intScaledP6N-uint(utils.Log10(int(fractional)))-leadingZeroes)))))
	if !negative {
		return a
	}
	return -a
}

// NewMFN creates a new FP object.
// The number will have at least minLeadingZeroes numbers after the dot.
//
// Arguments:
// - decimal - the number before the dot;
// - fractional - the number after the dot;
// - minFractionalNumbers - the minimum number of leading zeroes in the fractional part.
// - negative - if the number is negative.
//
// Examples:
// - NewMFN(2, 1, 2, false).String() = "2.01"
// - NewMFN(2, 10, 2, false).String() = "2.10"
// - NewMFN(2, 100, 2, false).String() = "2.100"
func NewMFN(decimal, fractional, minFractionalNumbers uint, negative bool) FP {
	var fractionalPow = uint(utils.Log10(int(fractional)))
	if fractionalPow < minFractionalNumbers {
		fractionalPow = minFractionalNumbers
	}

	var a = FP(int64(decimal)*intScaledP6Scale + (int64(fractional) * int64(utils.Pow10(int(intScaledP6N-fractionalPow)))))
	if !negative {
		return a
	}
	return -a
}

func NewFromFloat(f float64) FP {
	return FP(f * float64(intScaledP6Scale))
}

// Float64 converts the fixed point number to a float64.
func (i FP) Float64() float64 {
	if i.IsNaN() {
		return math.NaN()
	}

	return float64(i) / intScaledP6Scale
}

func (i FP) String() string {
	if i == 0 {
		return "0"
	}

	var neg = i < 0
	if neg {
		i = -i
	}

	var a = make([]byte, 21)
	var (
		n  int
		n1 = -1
	)
	for ; n < len(a) && (i > 0 || n < intScaledP6N+1); n++ {
		var b = byte(i % 10)
		var wasEmpty = n1 == -1
		i /= 10

		if b != 0 || n1 != -1 || n >= intScaledP6N {
			n1++
		}
		if n1 != -1 {
			if n == intScaledP6N && !wasEmpty {
				a[len(a)-n1-1] = '.'
				n1++
			}
			a[len(a)-n1-1] = '0' + b
		}
	}
	if neg {
		n1++
		a[len(a)-n1-1] = '-'
	}

	return string(a[len(a)-n1-1:])
}

/*
	Math operations
*/

func (i FP) Multiply(i2 FP) FP {
	var a1, b1 = i / intScaledP6Scale, i % intScaledP6Scale
	var a2, b2 = i2 / intScaledP6Scale, i2 % intScaledP6Scale

	return (a1 * b2) + (a2 * b1) + (a1*a2)*intScaledP6Scale + (b1*b2)/intScaledP6Scale
}

func (i FP) Divide(i2 FP) FP {
	return NewFromFloat(i.Float64() / i2.Float64())
}

// Floor cuts off any digits after n
func (i FP) Floor(n int) FP {
	if i < 0 {
		return -(-i).Ceil(n)
	}
	if n < 0 {
		n = 0
	}

	return i.FloorPositive(n)
}

// FloorPositive works correctly only on numbers [0, +inf)
func (i FP) FloorPositive(n int) FP {
	var pow = FP(utils.Pow10(intScaledP6N - n))
	var fractional = (i % intScaledP6Scale) / pow * pow

	return i/intScaledP6Scale*intScaledP6Scale + fractional
}

// Ceil rounds a number up to the next largest number
func (i FP) Ceil(n int) FP {
	if i < 0 {
		return -(-i).Floor(n)
	}

	if n < 0 {
		n = 0
	}

	return i.CeilPositive(n)
}

// CeilPositive works correctly only on numbers [0, +inf)
func (i FP) CeilPositive(n int) FP {
	var pow = FP(utils.Pow10(intScaledP6N - n))
	var fractional = (i % intScaledP6Scale) / pow * pow

	if i%pow != 0 {
		fractional += pow
	}
	return i/intScaledP6Scale*intScaledP6Scale + fractional
}

func (i FP) Abs() FP {
	if i < 0 {
		return -i
	}
	return i
}

/*
	NaN support
*/

func (i FP) IsNaN() bool {
	return i == NaN
}

func ParseIntScaledP6NaN(src string) FP {
	if src == "null" {
		return NaN
	}
	return Parse(src)
}

func (i FP) StringNaN() string {
	if i.IsNaN() {
		return "null"
	}
	return i.String()
}

func wrapMathWithNaN(i, i2 FP, f func() FP) FP {
	if i.IsNaN() || i2.IsNaN() {
		return NaN
	}
	return f()
}

func wrapMathWithNaNSingle(i FP, f func() FP) FP {
	if i.IsNaN() {
		return NaN
	}
	return f()
}

func (i FP) AddNaN(i2 FP) FP {
	return wrapMathWithNaN(i, i2, func() FP { return i + i2 })
}

func (i FP) SubtractNaN(i2 FP) FP {
	return wrapMathWithNaN(i, i2, func() FP { return i - i2 })
}

func (i FP) MultiplyNaN(i2 FP) FP {
	return wrapMathWithNaN(i, i2, func() FP { return i.Multiply(i2) })
}

func (i FP) DivideNaN(i2 FP) FP {
	return wrapMathWithNaN(i, i2, func() FP { return i.Divide(i2) })
}

func (i FP) FloorNaN(n int) FP {
	return wrapMathWithNaNSingle(i, func() FP { return i.Floor(n) })
}

func (i FP) CeilNaN(n int) FP {
	return wrapMathWithNaNSingle(i, func() FP { return i.Ceil(n) })
}

/*
	Interfaces implementations
*/

func (i *FP) UnmarshalJSON(bytes []byte) error {
	*i = Parse(string(bytes))
	return nil
}

func (i *FP) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	var rawBytes = in.Raw()
	*i = Parse(*(*string)(unsafe.Pointer(&rawBytes)))
	if isTopLevel {
		in.Consumed()
	}
}

func (i FP) MarshalJSON() ([]byte, error) {
	return []byte(i.String()), nil
}

func (i FP) MarshalEasyJSON(w *jwriter.Writer) {
	var a = i.String()
	w.Raw(*(*[]byte)(unsafe.Pointer(&a)), nil)
}
