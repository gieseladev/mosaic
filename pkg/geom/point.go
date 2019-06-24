package geom

import (
	"fmt"
	"math"
)

// A Point represents a point using cartesian coordinates.
type Point struct {
	X, Y float64
}

// NewPoint creates a new point from the given coordinates.
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// NewPointFromPolar creates a new point from the given polar coordinates.
func NewPointFromPolar(radius, angle float64) Point {
	return NewPoint(radius*math.Cos(angle), radius*math.Sin(angle))
}

func (p Point) String() string {
	return fmt.Sprintf("<%.3f, %.3f>", p.X, p.Y)
}

// XY returns the x and y coordinates.
func (p *Point) XY() (float64, float64) {
	return p.X, p.Y
}

// Polar returns the point in polar form.
// Return order: radius, angle
func (p *Point) Polar() (float64, float64) {
	x, y := p.XY()
	r := math.Sqrt(x*x + y*y)
	a := math.Atan2(y, x)
	return r, a
}

// Add returns a new point with the given point added.
func (p Point) Add(other Point) Point {
	return NewPoint(p.X+other.X, p.Y+other.Y)
}

// Sub returns a new point with the given point subtracted.
func (p Point) Sub(other Point) Point {
	return NewPoint(p.X-other.X, p.Y-other.Y)
}

// Mul returns a new point multiplied with the factor.
func (p Point) Mul(factor float64) Point {
	p.X *= factor
	p.Y *= factor
	return p
}

// Scale multiplies two points component-wise
func (p Point) Scale(other Point) Point {
	p.X *= other.X
	p.Y *= other.Y
	return p
}
