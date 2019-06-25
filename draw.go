package mosaic

import (
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic/pkg/geom"
)

// drawSlice draws a circle slice to the given context.
func drawSlice(dc *gg.Context, centerX, centerY, radius, angleStart, angleEnd float64) {
	dc.NewSubPath()
	dc.MoveTo(centerX, centerY)
	dc.DrawArc(centerX, centerY, radius, angleStart, angleEnd)
	dc.ClosePath()
}

// drawPolygon draws a polygon
func drawPolygon(dc *gg.Context, pg geom.Polygon) {
	if pg.Empty() {
		return
	}

	dc.NewSubPath()
	p := pg.Vertices[0]
	dc.MoveTo(p.XY())

	for _, p := range pg.Vertices[1:] {
		dc.LineTo(p.XY())
	}

	dc.ClosePath()
}
