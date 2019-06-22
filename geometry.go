package mosaic

import "math"

func minSide(width, height int) int {
	if width > height {
		return width
	} else {
		return height
	}
}

func innerSquareRadius(r float64) float64 {
	return r / 2
}

func outerSquareRadius(r float64) float64 {
	return math.Sqrt(2*(r*r)) / 2
}
