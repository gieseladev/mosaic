package mosaic

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"math"
)

type Composer interface {
	Compose(dc *gg.Context, images ...image.Image) error
}

type ComposerFunc func(dc *gg.Context, images ...image.Image) error

func (f ComposerFunc) Compose(dc *gg.Context, images ...image.Image) error {
	return f(dc, images...)
}

func CirclesPie(dc *gg.Context, images ...image.Image) error {
	w := dc.Width()
	h := dc.Height()

	s := minSide(w, h)

	angle := 2 * math.Pi / float64(len(images))

	maskDC := gg.NewContext(w, h)
	radius := innerSquareRadius(float64(s))
	midX, midY := float64(w)/2, float64(h)/2

	for i, img := range images {
		// TODO fill to different square
		img = imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)
		imgDC := gg.NewContextForImage(img)

		maskDC.Clear()
		startAngle := float64(i) * angle
		drawSlice(maskDC, midX, midY, radius, startAngle, startAngle+angle)
		maskDC.Fill()

		err := dc.SetMask(maskDC.AsMask())
		if err != nil {
			return err
		}

		dc.DrawImageAnchored(imgDC.Image(), w/2, h/2, .5, .5)
	}

	return nil
}

func TilesPerfect(dc *gg.Context, images ...image.Image) error {
	w := dc.Width()
	h := dc.Height()
	n := int(math.Sqrt(float64(len(images))))

	imgW := w / n
	imgH := h / n

	for i, img := range images {
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
