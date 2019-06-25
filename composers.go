package mosaic

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic/pkg/geom"
	"image"
	"math"
)

var (
	ErrInvalidImageCount = errors.New("invalid number of images")
)

var (
	circleCornerAngles = []float64{0, geom.HalfPi, math.Pi, 3 * geom.HalfPi}
	circleCornerPoints = []geom.Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
)

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

	nH, nV := geom.FindBalancedFactors(len(images))
	imgW := w / nH
	imgH := h / nV

	for i, img := range images {
		column := i % nH
		row := i / nH
		img = imaging.Fill(img, imgW, imgH, imaging.Center, imaging.Lanczos)
		dc.DrawImage(img, column*imgW, row*imgH)
	}

	return nil
}

func TilesFocused(dc *gg.Context, images ...image.Image) error {
	if len(images) < 2 {
		return ErrInvalidImageCount
	}

	w, h := dc.Width(), dc.Height()

	totalSize := geom.NewPoint(float64(w), float64(h))

	evenDiff := len(images) % 2
	evenImages := len(images) - evenDiff
	unevenImages := len(images) - (1 - evenDiff)

	horizontalRatio := float64(unevenImages-1) / float64(unevenImages+1)
	verticalRatio := float64(evenImages-2) / float64(evenImages)

	focusSize := totalSize.Scale(geom.NewPoint(
		horizontalRatio,
		verticalRatio,
	))
	focusX, focusY := int(focusSize.X), int(focusSize.Y)

	otherSize := totalSize.Sub(focusSize)
	otherX, otherY := int(otherSize.X), int(otherSize.Y)

	focusImg := imaging.Fill(images[0], focusX, focusY, imaging.Center, imaging.Lanczos)
	dc.DrawImage(focusImg, 0, h-focusY)

	trImg := imaging.Fill(images[1], otherX, otherY, imaging.Center, imaging.Lanczos)
	dc.DrawImage(trImg, focusX, 0)

	i := 1
	for imgI := 2; imgI < len(images); imgI += 2 {
		topImg := imaging.Fill(images[imgI], otherX, otherY, imaging.Center, imaging.Lanczos)
		dc.DrawImageAnchored(topImg, w-i*otherX, 0, 1, 0)

		rightI := imgI + 1
		if rightI < len(images) {
			rightImg := imaging.Fill(images[rightI], otherX, otherY, imaging.Center, imaging.Lanczos)
			dc.DrawImage(rightImg, focusX, i*otherY)
		}

		i++
	}

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

		ComposerInfo{
			Composer: ComposerFunc(StripesVertical),
			Id:       "stripes-vertical",
			Name:     "Vertical (Stripes)",
		},

		ComposerInfo{
			Composer: ComposerFunc(StripesVerticalMulti),
			Id:       "stripes-vertical-multi",
			Name:     "Vertical Multi (Stripes)",
		},
	)

	if err != nil {
		panic(fmt.Sprintf("couldn't register all built-in composers: %v", err))
	}
}
