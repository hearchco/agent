package rank

import (
	"math"
)

const magic64 = 0x5FE6EB50C7B537A9

func fastInvSqrt64(n float64) float64 {
	if n < 0 {
		return math.NaN()
	}
	n2, th := n*0.5, float64(1.5)
	b := math.Float64bits(n)
	b = magic64 - (b >> 1)
	f := math.Float64frombits(b)
	f *= th - (n2 * f * f)
	return f
}
