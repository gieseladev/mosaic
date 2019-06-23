package mosaic

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"testing"
)

func loadImage(path string) (image.Image, error) {
	return gg.LoadImage(path)
}

func loadImages(paths ...string) ([]image.Image, error) {
	images := make([]image.Image, len(paths))

	for i, path := range paths {
		img, err := loadImage(path)
		if err != nil {
			return images, err
		}

		images[i] = img
	}

	return images, nil
}

func loadInputImages(t *testing.T, names ...string) []image.Image {
	for i := 0; i < len(names); i++ {
		names[i] = fmt.Sprintf("test/data/input/%s", names[i])
	}

	images, err := loadImages(names...)

	require.NoError(t, err, fmt.Sprintf("couldn't load input images: %q", names))
	return images
}

func loadOutputImage(t *testing.T, name string) image.Image {
	img, err := loadImage(fmt.Sprintf("test/data/output/%s", name))
	require.NoError(t, err, fmt.Sprintf("couldn't load output image: %q", name))
	return img
}

func assertImagesEqual(t *testing.T, expected, actual image.Image, msgAndArgs ...interface{}) bool {
	if !expected.Bounds().Eq(actual.Bounds()) {
		return assert.Fail(t,
			fmt.Sprintf("Image sizes not equal:\n"+
				"expected: %s\n"+
				"actual  : %s", expected.Bounds(), expected.Bounds()),
			msgAndArgs...)
	}

	bounds := expected.Bounds()

	totalPixels := bounds.Dx() * bounds.Dy()
	var diffPixels uint64
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			eR, eG, eB, eA := expected.At(x, y).RGBA()
			aR, aG, aB, aA := actual.At(x, y).RGBA()

			if eR != aR || eG != aG || eB != aB || eA != aA {
				diffPixels++
			}
		}
	}

	diff := float64(diffPixels) / float64(totalPixels)
	if diff >= .05 {
		return assert.Fail(t, fmt.Sprintf("Images too different: %g%%", 100*diff))
	}

	return true
}

func TestPieChart(t *testing.T) {
	images := loadInputImages(t, "b-martinez-744134.jpg", "j-crop-764891.jpg", "j-han-456323.jpg")

	dc := gg.NewContext(250, 250)
	err := CirclesPie(dc, images...)
	require.NoError(t, err)

	actual := dc.Image()
	expected := loadOutputImage(t, "circles_pie-250x.png")

	assertImagesEqual(t, expected, actual)
}

func TestTilesPerfect(t *testing.T) {
	images := loadInputImages(t, "m-spiske-78531.jpg", "j-han-456323.jpg", "s-imbrock-487035.jpg", "s-erixon-753182.jpg")

	dc := gg.NewContext(250, 250)
	err := TilesPerfect(dc, images...)
	require.NoError(t, err)

	actual := dc.Image()
	expected := loadOutputImage(t, "tiles_perfect-250x.png")

	assertImagesEqual(t, expected, actual)
}

func TestStripesVertical(t *testing.T) {
	images := loadInputImages(t, "j-crop-764891.jpg", "s-imbrock-487035.jpg", "s-erixon-753182.jpg")

	dc := gg.NewContext(250, 250)
	err := StripesVertical(dc, images...)
	require.NoError(t, err)

	actual := dc.Image()
	expected := loadOutputImage(t, "stripes_vertical-250x.png")

	assertImagesEqual(t, expected, actual)
}
