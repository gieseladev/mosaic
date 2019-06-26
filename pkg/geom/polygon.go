package geom

// A Polygon represents any polygon.
type Polygon struct {
	Vertices []Point
}

// Poly creates a new polygon containing the given points.
func Poly(points ...Point) Polygon {
	return Polygon{Vertices: points}
}

// Empty checks whether the polygon contains no points
func (pg Polygon) Empty() bool {
	return len(pg.Vertices) == 0
}

// BoundingRect returns the bounding rectangle
func (pg Polygon) BoundingRect() Rectangle {
	return RectContainingPoints(pg.Vertices...)
}

// mapVertices returns a polygon with the given function applied to each
// vertex.
func (pg Polygon) mapVertices(f func(vertex Point) Point) Polygon {
	vertices := make([]Point, len(pg.Vertices))
	for i, v := range pg.Vertices {
		vertices[i] = f(v)
	}

	return Poly(vertices...)
}

func (pg Polygon) Center() Point {
	if pg.Empty() {
		return Point{}
	}

	sumX, sumY := pg.Vertices[0].XY()

	for _, v := range pg.Vertices[1:] {
		sumX += v.X
		sumY += v.Y
	}

	return Pt(sumX, sumY).Div(float64(len(pg.Vertices)))
}

// Translate moves the polygon by the given amount.
func (pg Polygon) Translate(p Point) Polygon {
	return pg.mapVertices(func(vertex Point) Point {
		return vertex.Add(p)
	})
}

// Scale scales the polygon from the origin.
func (pg Polygon) Scale(factor float64) Polygon {
	return pg.mapVertices(func(vertex Point) Point {
		return vertex.Mul(factor)
	})
}

// ScaleFrom scales the polygon from a given point.
func (pg Polygon) ScaleFrom(factor float64, origin Point) Polygon {
	return pg.
		Translate(origin.Neg()).
		Scale(factor).
		Translate(origin)
}

// ScaleFromCenter scales the polygon from the center.
func (pg Polygon) ScaleFromCenter(factor float64) Polygon {
	return pg.ScaleFrom(factor, pg.Center())
}
