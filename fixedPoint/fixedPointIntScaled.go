package fixedPoint

import (
	"github.com/k773/utils"
	"strconv"
)

type IntScaled int64

const intScaledN = 6
const intScaledScale = 10 * 10 * 10 * 10 * 10 * 10

func ParseFixedPointToIntScaled(src string) IntScaled {
	a, b, c, n := ParseFixedPointRaw(src, intScaledN)
	return NewIntScaled(a, b, c, n)
}

func NewIntScaled(decimal, fractional, fractionalPow uint, negative bool) IntScaled {
	var a = IntScaled(int64(decimal)*intScaledScale + (int64(fractional) * int64(utils.Pow10(int(intScaledN-uint(utils.Log10(int(fractional)))-fractionalPow)))))
	if !negative {
		return a
	}
	return -a
}

func NewIntScaledFromFloat64(f float64) IntScaled {
	return IntScaled(f * float64(intScaledScale))
}

// String: further optimizations are possible
func (i IntScaled) String() string {
	i2 := utils.If(i < 0, -i, i)
	var firstDigit = -1
	var fractional = make([]byte, intScaledN)
	for n := 0; n < intScaledN; n++ {
		var mod10 = byte(i2 % 10)
		if firstDigit != -1 || mod10 != 0 {
			if firstDigit == -1 {
				firstDigit = n
			}
			fractional[intScaledN-n-1] = '0' + mod10
		}
		i2 /= 10
	}

	var decimal = strconv.Itoa(int(i / intScaledScale))
	if firstDigit == -1 {
		return decimal
	}
	return decimal + "." + string(fractional[:intScaledN-firstDigit])
}

func (i IntScaled) Multiply(i2 IntScaled) IntScaled {
	var a1, b1 = i / intScaledScale, i % intScaledScale
	var a2, b2 = i2 / intScaledScale, i2 % intScaledScale

	var res IntScaled
	if a1 != 0 && a2 != 0 {
		res += (a1*a2)*intScaledScale + b1*a2
	}
	if b1 != 0 && b2 != 0 {
		res += (a1 * b2) + (b1*b2)/intScaledScale
	}

	return res
}

func (i IntScaled) Float64() float64 {
	return float64(i) / intScaledScale
}

func (i IntScaled) Divide(i2 IntScaled) IntScaled {
	return NewIntScaledFromFloat64(i.Float64() / i2.Float64())
}

// Floor cuts off any digits after n
func (i IntScaled) Floor(n int) IntScaled {
	if i < 0 {
		return -(-i).Ceil(n)
	}
	if n < 0 {
		n = 0
	}

	var pow = IntScaled(utils.Pow10(intScaledN - n))
	var fractional = (i % intScaledScale) / pow * pow

	return i/intScaledScale*intScaledScale + fractional
}

// Ceil rounds a number up to the next largest number
func (i IntScaled) Ceil(n int) IntScaled {
	if i < 0 {
		return -(-i).Floor(n)
	}

	if n < 0 {
		n = 0
	}

	var pow = IntScaled(utils.Pow10(intScaledN - n))
	var fractional = (i % intScaledScale) / pow * pow

	if i%pow != 0 {
		fractional += pow
	}
	return i/intScaledScale*intScaledScale + fractional
}

func (i *IntScaled) UnmarshalJSON(bytes []byte) error {
	*i = ParseFixedPointToIntScaled(string(bytes))
	return nil
}

func (i IntScaled) MarshalJSON() ([]byte, error) {
	return []byte(i.String()), nil
}
