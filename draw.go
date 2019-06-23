package mosaic

import "github.com/fogleman/gg"

// drawSlice draws a circle slice to the given context.
func drawSlice(dc *gg.Context, centerX, centerY, radius, angleStart, angleEnd float64) {
	dc.MoveTo(centerX, centerY)
	dc.DrawArc(centerX, centerY, radius, angleStart, angleEnd)
	dc.LineTo(centerX, centerY)
}
