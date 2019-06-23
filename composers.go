package mosaic

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic/pkg/geom"
	"image"
	"math"
)

var circleCornerAngles = []float64{0, geom.HalfPi, math.Pi, 3 * geom.HalfPi}
var circleCornerPoints = []geom.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func CirclesPie(dc *gg.Context, images ...image.Image) error {
	w := dc.Width()
	h := dc.Height()

	var s int
	if w > h {
		s = h
	} else {
		s = w
	}

	angle := geom.TwoPi / float64(len(images))

	maskDC := gg.NewContext(w, h)
	radius := geom.InnerSquareRadius(float64(s))
	centerPoint := geom.NewPoint(float64(w)/2, float64(h)/2)

	for i, img := range images {
		maskDC.Clear()
		startAngle := float64(i) * angle
		endAngle := startAngle + angle
		drawSlice(maskDC, centerPoint.X, centerPoint.Y, radius, startAngle, endAngle)
		maskDC.Fill()

		err := dc.SetMask(maskDC.AsMask())
		if err != nil {
			return err
		}

		rect := geom.NewRectContainingPoints(
			centerPoint,
			geom.NewPointFromPolar(radius, startAngle).Add(centerPoint),
			geom.NewPointFromPolar(radius, endAngle).Add(centerPoint),
		)

		// ensure the rect covers all of the circle
		for i, angle := range circleCornerAngles {
			if geom.AngleStrictlyBetween(angle, startAngle, endAngle) {
				pp := circleCornerPoints[i].Mul(radius).Add(centerPoint)
				rect = rect.GrowToContain(pp)
			}
		}

		img = imaging.Fill(img, int(rect.DX()), int(rect.DY()), imaging.Center, imaging.Lanczos)
		dc.DrawImage(img, int(rect.Min.X), int(rect.Min.Y))
	}

	return nil
}

func TilesPerfect(dc *gg.Context, images ...image.Image) error {
	w := dc.Width()
	h := dc.Height()
	n := int(math.Sqrt(float64(len(images))))

	imgW := w / n
	imgH := h / n

	maxI := n * n
	for i, img := range images {
		if i > maxI {
			break
		}

		column := i % n
		row := i / n
		img = imaging.Fill(img, imgW, imgH, imaging.Center, imaging.Lanczos)
		dc.DrawImage(img, column*imgW, row*imgH)
	}

	return nil
}

func TilesFocused(dc *gg.Context, images ...image.Image) error {
	// TODO
	return nil
}

func TilesDiamond(dc *gg.Context, images ...image.Image) error {
	// TODO
	return nil
}

func StripesVertical(dc *gg.Context, images ...image.Image) error {
	w := dc.Width()
	h := dc.Height()
	stripeWidth := float64(w) / float64(len(images))

	maskDC := gg.NewContext(w, h)

	for i, img := range images {
		img = imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)

		iF64 := float64(i)
		maskDC.Clear()
		maskDC.DrawRectangle(iF64*stripeWidth, 0, stripeWidth, float64(h))
		maskDC.Fill()

		err := dc.SetMask(maskDC.AsMask())
		if err != nil {
			return err
		}

		dc.DrawImage(img, 0, 0)
	}

	return nil
}

func StripesVerticalMulti(dc *gg.Context, images ...image.Image) error {
	// TODO
	return nil
}

func init() {
	err := RegisterComposer(
		ComposerInfo{
			Composer: ComposerFunc(CirclesPie),
			Id:       "circles-pie",
			Name:     "Pie (Circle)",
		},

		ComposerInfo{
			Composer: ComposerFunc(TilesPerfect),
			Id:       "tiles-perfect",
			Name:     "Perfect (Tile)",
		},
		ComposerInfo{
			Composer: ComposerFunc(TilesFocused),
			Id:       "tiles-focused",
			Name:     "Focused (Tile)",
		},
		ComposerInfo{
			Composer: ComposerFunc(TilesDiamond),
			Id:       "tiles-diamond",
			Name:     "Diamond (Tile)",
		},
	)

	if err != nil {
		panic(fmt.Sprintf("couldn't register all built-in composers: %v", err))
	}
}
