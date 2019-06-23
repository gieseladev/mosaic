package geom

import (
	"fmt"
	"math"
)

// A Rectangle represents a rectangle using 2 corner points.
type Rectangle struct {
	Min, Max Point
}

// NewRectFromSize creates a new rectangle with the given size.
func NewRectFromSize(width, height float64) Rectangle {
	return Rectangle{
		Max: Point{width, height},
	}
}

// NewRectContainingPoints finds the smallest rectangle containing all given
// points.
func NewRectContainingPoints(points ...Point) Rectangle {
	if len(points) == 0 {
		return Rectangle{}
	}

	p, points := points[len(points)-1], points[:len(points)-1]
	return Rectangle{Min: p, Max: p}.GrowToContain(points...)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%v / %v)", r.Min, r.Max)
}

// DX returns the difference between the x coordinates (i.e. the width).
func (r Rectangle) DX() float64 {
	return r.Max.X - r.Min.X
}

// DY returns the difference between the y coordinates (i.e. the height).
func (r Rectangle) DY() float64 {
	return r.Max.Y - r.Min.Y
}

// MinSide return the length of the smaller side.
func (r Rectangle) MinSide() float64 {
	return math.Min(r.DX(), r.DY())
}

// MaxSide returns the length of the bigger side.
func (r Rectangle) MaxSide() float64 {
	return math.Max(r.DX(), r.DY())
}

// GrowToContain returns a new rectangle expanded to contain the given points.
// If the rectangle already contains the points, a copy is returned.
func (r Rectangle) GrowToContain(points ...Point) Rectangle {
	for _, p := range points {
		x, y := p.XY()

		if x < r.Min.X {
			r.Min.X = x
		}
		if y < r.Min.Y {
			r.Min.Y = y
		}

		if x > r.Max.X {
			r.Max.X = x
		}
		if y > r.Max.Y {
			r.Max.Y = y
		}
	}

	return r
}
