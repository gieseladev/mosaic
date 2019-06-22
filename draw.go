package mosaic

import "github.com/fogleman/gg"

func drawSlice(dc *gg.Context, centerX, centerY, radius, angleStart, angleEnd float64) {
	dc.MoveTo(centerX, centerY)
	dc.DrawArc(centerX, centerY, radius, angleStart, angleEnd)
	dc.LineTo(centerX, centerY)
}
