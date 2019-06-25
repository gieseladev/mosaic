package geom

// A Polygon represents any shape by connecting the points.
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

func (pg Polygon) mapVertices(f func(vertex Point) Point) []Point {
	vertices := make([]Point, len(pg.Vertices))
	for i, v := range pg.Vertices {
		vertices[i] = f(v)
	}

	return vertices
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
	return Poly(pg.mapVertices(func(vertex Point) Point {
		return vertex.Add(p)
	})...)
}

func (pg Polygon) Scale(factor float64) Polygon {
	return Poly(pg.mapVertices(func(vertex Point) Point {
		return vertex.Mul(factor)
	})...)
}

func (pg Polygon) ScaleAround(factor float64, origin Point) Polygon {
	return pg.
		Translate(origin.Neg()).
		Scale(factor).
		Translate(origin)
}

func (pg Polygon) ScaleAroundCenter(factor float64) Polygon {
	return pg.ScaleAround(factor, pg.Center())
}
