package fixedPoint

func ParseFixedPointRaw(s string, maxFractionalPrecision int) (decimal, fractional, zeroesBeforeFractional uint, negative bool) {
	var n = 0
	var swapDot bool
f:
	for i := 0; i < len(s); i++ {
		switch c := s[i]; true {
		case c >= '0' && c <= '9':
			if swapDot {
				// Limit max fractional part
				if n == maxFractionalPrecision {
					break f
				}
				n++

				// Counting number of zeroes before the fractional part start
				if c == '0' && fractional == 0 {
					zeroesBeforeFractional++
				}
				fractional = fractional*10 + uint(c-'0')
			} else {
				decimal = decimal*10 + uint(c-'0')
			}
		case c == '.':
			if !swapDot {
				swapDot = true
			}
		case c == '-':
			if !negative {
				negative = true
			}
		}
	}

	return
}
