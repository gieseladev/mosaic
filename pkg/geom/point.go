package geom

import (
	"fmt"
	"math"
)

// A Point represents a point using cartesian coordinates.
type Point struct {
	X, Y float64
}

// Pt creates a new point from the given coordinates.
func Pt(x, y float64) Point {
	return Point{X: x, Y: y}
}

// PtFromPolar creates a new point from the given polar coordinates.
func PtFromPolar(radius, angle float64) Point {
	return Pt(radius*math.Cos(angle), radius*math.Sin(angle))
}

func (p Point) String() string {
	return fmt.Sprintf("<%.3f, %.3f>", p.X, p.Y)
}

// XY returns the x and y coordinates.
func (p Point) XY() (float64, float64) {
	return p.X, p.Y
}

// Neg returns the negative point
func (p Point) Neg() Point {
	return Pt(-p.X, -p.Y)
}

// Polar returns the point in polar form.
// Return order: radius, angle
func (p Point) Polar() (float64, float64) {
	x, y := p.XY()
	r := math.Sqrt(x*x + y*y)
	a := math.Atan2(y, x)
	return r, a
}

// Add returns a new point with the given point added.
func (p Point) Add(other Point) Point {
	return Pt(p.X+other.X, p.Y+other.Y)
}

// Sub returns a new point with the given point subtracted.
func (p Point) Sub(other Point) Point {
	return Pt(p.X-other.X, p.Y-other.Y)
}

// Mul returns a new point multiplied with the factor.
func (p Point) Mul(factor float64) Point {
	p.X *= factor
	p.Y *= factor
	return p
}

// Div returns a new point divided by the divisor.
func (p Point) Div(divisor float64) Point {
	p.X /= divisor
	p.Y /= divisor
	return p
}

// Scale multiplies two points component-wise
func (p Point) Scale(other Point) Point {
	p.X *= other.X
	p.Y *= other.Y
	return p
}

// Rotate rotates a point counterclockwise around the origin
func (p Point) Rotate(angle float64) Point {
	sin, cos := math.Sincos(angle)

	return Pt(p.X*cos-p.Y*sin, p.X*sin+p.Y*cos)
}

// RotateAround rotates a point counterclockwise around another point.
func (p Point) RotateAround(angle float64, origin Point) Point {
	return p.
		Sub(origin).
		Rotate(angle).
		Add(origin)
}
