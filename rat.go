package utils

import (
	"golang.org/x/exp/constraints"
	"math"
	"strconv"
)

/*
	If max speed of log10 is needed for small numbers, increase array size to approximately maximum number that will be calculated
*/

var fastLog10 = [0]int{}

func init() {
	for i := range fastLog10 {
		fastLog10[i] = Log10Loop(i)
	}
}

func ParseRat(s string, maxFractionalPrecision int) *Rat {
	var a, n = 0, 0
	var swapdot bool
	var negative bool
f:
	for i := 0; i < len(s); i++ {
		switch c := s[i]; true {
		case c >= '0' && c <= '9':
			if swapdot {
				if n == maxFractionalPrecision {
					break f
				}
				n++
			}
			a = a*10 + int(c-'0')
		case c == '.':
			if !swapdot {
				swapdot = true
			}
		case c == '-':
			if !negative {
				negative = true
			}
		}
	}
	//if fractionalPrecision > 0 && dn != fractionalPrecision {
	//	b *= Pow10(fractionalPrecision - dn)
	//}

	if negative {
		a = -a
	}
	return NewRat(a, n)
}

func ParseRat2(s string, maxFractionalPrecision int) (a int, b int) {
	var swapdot bool
	var dn int
f:
	for i := 0; i < len(s); i++ {
		switch c := s[i]; true {
		case c >= '0' && c <= '9':
			if swapdot {
				if dn == maxFractionalPrecision {
					break f
				}
				dn++
				b = b*10 + int(c-'0')
			} else {
				a = a*10 + int(c-'0')
			}
		case c == '.':
			if !swapdot {
				swapdot = true
			}
		}
	}
	//if fractionalPrecision > 0 && dn != fractionalPrecision {
	//	b *= Pow10(fractionalPrecision - dn)
	//}

	return
}

func RatToString(a, b int) string {
	return strconv.Itoa(a) + "." + strconv.Itoa(b)
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

var pow10 = []int{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e15, 1e16, 1e17, 1e18}

func Pow10(a int) int {
	if a < len(pow10) {
		return pow10[a]
	}
	return math.MaxInt64
}

func Log10(a int) int {
	if a < 0 {
		return -Log10(-a)
	}
	if a < len(fastLog10) {
		return fastLog10[a]
	}
	return Log10Loop(a)
}

func Log10Loop(a int) int {
	if a >= 1e18 {
		return 19
	}
	x, count := 10, 1
	for x <= a {
		x *= 10
		count++
	}
	return count
}

var tab64 = [64]int{
	63, 0, 58, 1, 59, 47, 53, 2,
	60, 39, 48, 27, 54, 33, 42, 3,
	61, 51, 37, 40, 49, 18, 28, 20,
	55, 30, 34, 11, 43, 14, 22, 4,
	62, 57, 46, 52, 38, 26, 32, 41,
	50, 36, 17, 19, 29, 10, 13, 21,
	56, 45, 25, 31, 35, 16, 9, 12,
	44, 24, 15, 8, 23, 7, 6, 5}

func Log2(value uint64) int {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32
	return tab64[((value-(value>>1))*0x07EDD5E59A4E28C2)>>58]
}

func BoolToNumber[T int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint | float32 | float64](val bool) T {
	var a T
	if val {
		a = 1
	}
	return a
}

type Rat struct {
	A, N         int
	DynamicPoint bool
}

func NewRat(a, n int) *Rat {
	return &Rat{A: a, N: n, DynamicPoint: false}
}

/*
	Tools
*/

func (r *Rat) Clone() *Rat {
	var r2 = *r
	return &r2
}

func (r *Rat) GrowTo(pow int) *Rat {
	if pow > 17 {
		pow = 17
	}
	p := pow - Log10(r.A)
	if p > 0 {
		r.A *= Pow10(p)
		r.N += p
	}
	return r
}

func (r *Rat) ShrinkTo(pow int) *Rat {
	if r.N > pow {
		d := r.N - pow
		r.A /= Pow10(d)
		r.N -= d
	}
	return r
}

func (r *Rat) ToPow(pow int) *Rat {
	if r.N > pow {
		return r.ShrinkTo(pow)
	} else if r.N < pow {
		return r.GrowTo(pow)
	}
	return r
}

// DynShrink shrinks until max steps are performed or last digit is non-zero
func (r *Rat) DynShrink(max int) {
	for i := 0; i < max && r.A%10 == 0; i++ {
		r.A /= 10
		r.N--
	}
}

/*
	Multiply
*/

func (r *Rat) Multiply(r2 *Rat) *Rat {
	r = r.Clone()
	if !r.DynamicPoint {
		r.A = (r.A * r2.A) / Pow10(r2.N)
	} else {
		r.N += r2.N
		r.A *= r2.A
		r.DynShrink(r2.N)
	}
	return r
}

/*
	Divide
*/

func (r *Rat) Divide(r2 *Rat) *Rat {
	r = r.Clone()
	if !r.DynamicPoint {
		a := Abs(r.N - (r.N - r2.N))
		r.A = (r.A * Pow10(a)) / r2.A
		r.N += a - r2.N
	} else {
		r.GrowTo(17)
		r.A /= r2.A
		r.N -= r2.N
	}
	return r
}

/*
	Add
*/

func (r *Rat) Add(r2 *Rat) *Rat {
	r = r.Clone()
	d := r.N - r2.N

	if d > 0 {
		r.A = (r2.A * Pow10(d)) + r.A
	} else if d < 0 {
		if r.DynamicPoint {
			r.A = (r.A * Pow10(-d)) + r2.A
			r.N += -d
		} else {
			r.A = ((r.A * Pow10(-d)) + r2.A) / Pow10(-d)
		}
	} else {
		r.A += r2.A
	}
	return r
}

func (r *Rat) AddToThis(r2 *Rat) {
	d := r.N - r2.N

	if d > 0 {
		r.A = (r2.A * Pow10(d)) + r.A
	} else if d < 0 {
		if r.DynamicPoint {
			r.A = (r.A * Pow10(-d)) + r2.A
			r.N += -d
		} else {
			r.A = ((r.A * Pow10(-d)) + r2.A) / Pow10(-d)
		}
	} else {
		r.A += r2.A
	}
}

func (r *Rat) AddToThisIfNotNegative(r2 *Rat) {
	if !r2.Negative() {
		r.AddToThis(r2)
	}
}

/*
	Sub
*/

func (r *Rat) Sub(r2 *Rat) *Rat {
	r = r.Clone()
	d := r.N - r2.N

	if d > 0 {
		r.A = r.A - (r2.A * Pow10(d))
	} else if d < 0 {
		if r.DynamicPoint {
			r.A = (r.A * Pow10(-d)) - r2.A
			r.N += -d
		} else {
			r.A = ((r.A * Pow10(-d)) - r2.A) / Pow10(-d)
		}
	} else {
		r.A -= r2.A
	}
	return r
}

// Compare returns 1 if r>r2; 0 if r==r2; -1 if r<r2
func (r *Rat) Compare(r2 *Rat) int {
	d := r.N - r2.N
	var n0, n1 int

	if d > 0 {
		n0, n1 = r.A, r2.A*Pow10(d)
	} else if d < 0 {
		n0, n1 = r.A*Pow10(-d), r2.A
	} else {
		n0, n1 = r.A, r2.A
	}

	if n0 > n1 {
		return 1
	} else if n0 == n1 {
		return 0
	} else {
		return -1
	}
}

func (r *Rat) Equals(r2 *Rat) bool {
	return r.Compare(r2) == 0
}

func (r *Rat) Less(r2 *Rat) bool {
	return r.Compare(r2) == -1
}

func (r *Rat) LessOrEquals(r2 *Rat) bool {
	a := r.Compare(r2)
	return a == -1 || a == 0
}

func (r *Rat) Greater(r2 *Rat) bool {
	return r.Compare(r2) == 1
}

func (r *Rat) GreaterOrEquals(r2 *Rat) bool {
	a := r.Compare(r2)
	return a == 1 || a == 0
}

func (r *Rat) Parts() (a, b int) {
	n := Pow10(r.N)
	return r.A / n, Abs(r.A % n)
}

func (r *Rat) Number() int {
	return r.A
}

func (r *Rat) Float64() float64 {
	n := Pow10(r.N)
	return float64(r.A/n) + float64(r.A%n)/float64(n)
}

func (r *Rat) Negative() bool {
	return r.A < 0
}

func (r *Rat) String() string {
	n := Pow10(r.N)
	if r.N != 0 {
		return strconv.Itoa(r.A/n) + "." + PadLeft(strconv.Itoa(Abs(r.A%n)), r.N, '0')
	}
	return strconv.Itoa(r.A / n)
}

func (r *Rat) MarshalJSON() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Rat) UnmarshalJSON(a []byte) error {
	b := ParseRat(string(a), -1)
	r.A, r.N = b.A, b.N
	return nil
}
