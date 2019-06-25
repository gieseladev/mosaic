package geom

import (
	"fmt"
	"math"
)

// A Rectangle represents a rectangle using 2 corner points.
type Rectangle struct {
	Min, Max Point
}

// RectWithSideLengths creates a new rectangle with the given size.
func RectWithSideLengths(p Point) Rectangle {
	return Rectangle{Max: p}
}

// SquareWithSideLen creates a new square with the given size.
func SquareWithSideLen(side float64) Rectangle {
	return RectWithSideLengths(Pt(side, side))
}

// RectContainingPoints finds the smallest rectangle containing all given
// points.
func RectContainingPoints(points ...Point) Rectangle {
	if len(points) == 0 {
		return Rectangle{}
	}

	return Rectangle{Min: points[0], Max: points[0]}.
		GrowToContain(points[1:]...)
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

// TopLeft returns the top left corner point
func (r Rectangle) TopLeft() Point {
	return r.Min
}

// TopLeft returns the top right corner point
func (r Rectangle) TopRight() Point {
	return Pt(r.Max.X, r.Min.Y)
}

// TopLeft returns the bottom right corner point
func (r Rectangle) BottomRight() Point {
	return r.Max
}

// TopLeft returns the bottom left corner point
func (r Rectangle) BottomLeft() Point {
	return Pt(r.Min.X, r.Max.Y)
}

// Vertices returns a slice containing all four corner points
// The order is clockwise starting with the top left corner.
func (r Rectangle) Vertices() []Point {
	return []Point{r.TopLeft(), r.TopRight(), r.BottomRight(), r.BottomLeft()}
}

// Center returns the center of the rectangle.
func (r Rectangle) Center() Point {
	return r.Min.
		Add(r.Max).
		Mul(.5)
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

// Translate moves the rectangle around by the given point.
func (r Rectangle) Translate(p Point) Rectangle {
	return Rectangle{
		Min: r.Min.Add(p),
		Max: r.Max.Add(p),
	}
}

// Scale scales the rectangle.
func (r Rectangle) Scale(factor float64) Rectangle {
	return Rectangle{
		Min: r.Min.Mul(factor),
		Max: r.Max.Mul(factor),
	}
}

// ScaleCenter scales the rectangle from the center.
func (r Rectangle) ScaleCenter(factor float64) Rectangle {
	center := r.Center()

	return r.
		Translate(center.Neg()).
		Scale(factor).
		Translate(center)

}

// RotateAround returns the four vertices of the rectangle after a
// counterclockwise rotation around the given point.
func (r Rectangle) RotateAround(angle float64, origin Point) Polygon {
	return Poly(
		r.TopLeft().RotateAround(angle, origin),
		r.TopRight().RotateAround(angle, origin),
		r.BottomRight().RotateAround(angle, origin),
		r.BottomLeft().RotateAround(angle, origin),
	)
}

// RotateAroundCenter returns the four vertices of the rectangle after a
// counterclockwise rotation around the center.
func (r Rectangle) RotateAroundCenter(angle float64) Polygon {
	return r.RotateAround(angle, r.Center())
}

// InnerSquare returns a new square which fits inside of the rectangle.
func (r Rectangle) InnerCenterSquare() Rectangle {
	w, h := r.DX(), r.DY()
	s := w
	if w > h {
		s = h
	}

	halfS := float64(s) / 2
	diag := Pt(halfS, halfS)

	c := r.Center()

	return Rectangle{
		Min: c.Sub(diag),
		Max: c.Add(diag),
	}
}
