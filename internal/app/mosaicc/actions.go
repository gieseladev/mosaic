package mosaicc

import (
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic"
	"image"
)

func GenerateComposerShowcase(images []image.Image) (image.Image, error) {
	composers := mosaic.GetComposers()

	for _, composer := range composers {
		dc := gg.NewContext(100, 100)
		err := composer.Compose(dc, images...)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
