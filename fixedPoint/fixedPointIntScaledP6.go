/*
	IntScaledP6 implements fixed point math with precision up to 6 decimal places after the dot
*/

package fixedPoint

import (
	"github.com/k773/utils"
	"math"
)

type IntScaledP6 int64

// Constants
// If you want to change those, duplicate the type.
const intScaledP6N = 6
const intScaledP6Scale = 10 * 10 * 10 * 10 * 10 * 10
const NaN = -(1<<63 - 1) // bits set are identical to uint64(1<<64-1)

func ParseIntScaledP6(src string) IntScaledP6 {
	a, b, c, n := ParseFixedPointRaw(src, intScaledP6N)
	return NewIntScaledP6(a, b, c, n)
}

// NewIntScaledP6 creates a new IntScaledP6 object; fractionalPow - the number of leading zeroes in the fractional part
func NewIntScaledP6(decimal, fractional, fractionalPow uint, negative bool) IntScaledP6 {
	var a = IntScaledP6(int64(decimal)*intScaledP6Scale + (int64(fractional) * int64(utils.Pow10(int(intScaledP6N-uint(utils.Log10(int(fractional)))-fractionalPow)))))
	if !negative {
		return a
	}
	return -a
}

func NewIntScaledP6FromFloat64(f float64) IntScaledP6 {
	return IntScaledP6(f * float64(intScaledP6Scale))
}

func (i IntScaledP6) Float64() float64 {
	if i.IsNaN() {
		return math.NaN()
	}

	return float64(i) / intScaledP6Scale
}

// String: further optimizations are possible
func (i IntScaledP6) String() string {
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

func (i IntScaledP6) Multiply(i2 IntScaledP6) IntScaledP6 {
	var a1, b1 = i / intScaledP6Scale, i % intScaledP6Scale
	var a2, b2 = i2 / intScaledP6Scale, i2 % intScaledP6Scale

	return (a1 * b2) + (a2 * b1) + (a1*a2)*intScaledP6Scale + (b1*b2)/intScaledP6Scale
}

func (i IntScaledP6) Divide(i2 IntScaledP6) IntScaledP6 {
	return NewIntScaledP6FromFloat64(i.Float64() / i2.Float64())
}

// Floor cuts off any digits after n
func (i IntScaledP6) Floor(n int) IntScaledP6 {
	if i < 0 {
		return -(-i).Ceil(n)
	}
	if n < 0 {
		n = 0
	}

	var pow = IntScaledP6(utils.Pow10(intScaledP6N - n))
	var fractional = (i % intScaledP6Scale) / pow * pow

	return i/intScaledP6Scale*intScaledP6Scale + fractional
}

// Ceil rounds a number up to the next largest number
func (i IntScaledP6) Ceil(n int) IntScaledP6 {
	if i < 0 {
		return -(-i).Floor(n)
	}

	if n < 0 {
		n = 0
	}

	var pow = IntScaledP6(utils.Pow10(intScaledP6N - n))
	var fractional = (i % intScaledP6Scale) / pow * pow

	if i%pow != 0 {
		fractional += pow
	}
	return i/intScaledP6Scale*intScaledP6Scale + fractional
}

func (i IntScaledP6) Abs() IntScaledP6 {
	if i < 0 {
		return -i
	}
	return i
}

/*
	NaN support
*/

func (i IntScaledP6) IsNaN() bool {
	return i == NaN
}

func ParseIntScaledP6NaN(src string) IntScaledP6 {
	if src == "null" {
		return NaN
	}
	return ParseIntScaledP6(src)
}

func (i IntScaledP6) StringNaN() string {
	if i.IsNaN() {
		return "null"
	}
	return i.String()
}

func wrapMathWithNaN(i, i2 IntScaledP6, f func() IntScaledP6) IntScaledP6 {
	if i.IsNaN() || i2.IsNaN() {
		return NaN
	}
	return f()
}

func wrapMathWithNaNSingle(i IntScaledP6, f func() IntScaledP6) IntScaledP6 {
	if i.IsNaN() {
		return NaN
	}
	return f()
}

func (i IntScaledP6) AddNaN(i2 IntScaledP6) IntScaledP6 {
	return wrapMathWithNaN(i, i2, func() IntScaledP6 { return i + i2 })
}

func (i IntScaledP6) SubtractNaN(i2 IntScaledP6) IntScaledP6 {
	return wrapMathWithNaN(i, i2, func() IntScaledP6 { return i - i2 })
}

func (i IntScaledP6) MultiplyNaN(i2 IntScaledP6) IntScaledP6 {
	return wrapMathWithNaN(i, i2, func() IntScaledP6 { return i.Multiply(i2) })
}

func (i IntScaledP6) DivideNaN(i2 IntScaledP6) IntScaledP6 {
	return wrapMathWithNaN(i, i2, func() IntScaledP6 { return i.Divide(i2) })
}

func (i IntScaledP6) FloorNaN(n int) IntScaledP6 {
	return wrapMathWithNaNSingle(i, func() IntScaledP6 { return i.Floor(n) })
}

func (i IntScaledP6) CeilNaN(n int) IntScaledP6 {
	return wrapMathWithNaNSingle(i, func() IntScaledP6 { return i.Ceil(n) })
}

/*
	Interfaces implementations
*/

func (i *IntScaledP6) UnmarshalJSON(bytes []byte) error {
	*i = ParseIntScaledP6NaN(string(bytes))
	return nil
}

func (i IntScaledP6) MarshalJSON() ([]byte, error) {
	return []byte(i.StringNaN()), nil
}
