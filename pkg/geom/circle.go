package geom

import "math"

// InnerSquareRadius returns the inner radius of a square.
func InnerSquareRadius(r float64) float64 {
	return r / 2
}

// OuterSquareRadius returns the outer radius of a square.
func OuterSquareRadius(r float64) float64 {
	return math.Sqrt(2*(r*r)) / 2
}

// QuarterPi represents PI / 4
const QuarterPi = math.Pi / 4

// HalfPi represents PI / 2
const HalfPi = math.Pi / 2

// TwoPi represents 2 * PI
const TwoPi = 2 * math.Pi

// AngleStrictlyBetween checks whether the given angle is strictly between
// two other angles.
func AngleStrictlyBetween(angle, low, high float64) bool {
	highRel := high - low
	if highRel < 0 {
		highRel += TwoPi
	}

	angleRel := angle - low
	if angleRel < 0 {
		angleRel += TwoPi
	}

	return angleRel < highRel
}
