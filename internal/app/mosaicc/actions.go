package mosaicc

import (
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic"
	"github.com/gieseladev/mosaic/pkg/geom"
	"image"
	"image/color"
)

// GenerateComposerShowcase generates an image containing a sample of all
// composers.
func GenerateComposerShowcase(images []image.Image) (image.Image, error) {
	composers := mosaic.GetComposers()
	panelsX, panelsY := geom.FindBalancedFactors(len(composers))

	panelWidth, panelHeight := 200, 300
	marginX, marginY := 15, 15

	fontPoints := float64(panelWidth) / 6
	lineSpacing := 1.05

	dc := gg.NewContext(panelsX*panelWidth+(panelsX-1)*marginX,
		panelsY*panelHeight+(panelsY-1)*marginY)

	for i, composer := range composers {
		compImages := images[:composer.RecommendImageCount(len(images))]

		compositionDC := gg.NewContext(panelWidth, panelWidth)
		err := composer.Compose(compositionDC, compImages...)
		if err != nil {
			return nil, err
		}

		composerDC := gg.NewContext(panelWidth, panelHeight)

		composerDC.DrawImage(compositionDC.Image(), 0, panelHeight-panelWidth)

		// TODO nonononono, this ain't it, chief
		if err = composerDC.LoadFontFace("/windows/fonts/arial.ttf", fontPoints); err != nil {
			return nil, err
		}

		composerDC.SetColor(color.Black)
		composerDC.DrawStringWrapped(composer.Name, 0, float64(panelHeight-panelWidth-marginY), 0, 1, float64(panelWidth), lineSpacing, gg.AlignCenter)

		column := i % panelsX
		row := i / panelsX

		posX := column * (panelWidth + marginX)
		posY := row * (panelHeight + marginY)

		dc.DrawImage(composerDC.Image(), posX, posY)
	}

	return dc.Image(), nil
}
