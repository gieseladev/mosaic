package geom

import "math"

// FindBalancedFactors finds the two factors a and b of the given number
// such that the difference between a and b is minimal
// (i.e. a = b = sqrt(value) for square numbers)
// The return order is a, b such that a <= b
func FindBalancedFactors(value int) (int, int) {
	a := int(math.Ceil(math.Sqrt(float64(value))))

	for {
		b := value / a

		if a*b == value {
			return a, b
		}

		a -= 1
	}
}
