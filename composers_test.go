package mosaic

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
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

func loadInputImages(t *testing.T, names ...string) ([]image.Image, bool) {
	for i := 0; i < len(names); i++ {
		names[i] = fmt.Sprintf("test/data/input/%s", names[i])
	}

	images, err := loadImages(names...)

	ok := assert.NoError(t, err, fmt.Sprintf("couldn't load input images: %v", names))
	return images, ok
}

func loadOutputImage(t *testing.T, name string) (image.Image, bool) {
	img, err := loadImage(fmt.Sprintf("test/data/output/%s", name))
	ok := assert.NoError(t, err, fmt.Sprintf("couldn't load output image: %q", name))
	return img, ok
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

type ComposerTest struct {
	ComposerID string
	testName   string

	InputImageNames []string
	contextWidth    int
	contextHeight   int

	outputImageName string
}

func (c *ComposerTest) ContextWidth() int {
	if c.contextWidth == 0 {
		c.contextWidth = 250
	}

	return c.contextWidth
}

func (c *ComposerTest) ContextHeight() int {
	if c.contextHeight == 0 {
		c.contextHeight = 250
	}

	return c.contextHeight
}

func (c *ComposerTest) TestName() string {
	if c.testName == "" {
		c.testName = fmt.Sprintf("%s-%d-%dx%d",
			c.ComposerID,
			len(c.InputImageNames),
			c.ContextWidth(), c.ContextHeight(),
		)
	}

	return c.testName
}

func (c *ComposerTest) ExpectedImage(t *testing.T) (image.Image, bool) {
	if c.outputImageName == "" {
		c.outputImageName = fmt.Sprintf("%s.png", c.TestName())
	}

	return loadOutputImage(t, c.outputImageName)
}

func (c *ComposerTest) saveActualImage(t *testing.T, dc *gg.Context) bool {
	err := dc.SavePNG(fmt.Sprintf("%s/%s-actual.png",
		"test/data/output",
		c.TestName(),
	))

	return assert.NoError(t, err, "couldn't save actual output")
}

func (c *ComposerTest) Test(t *testing.T) (ok bool) {
	composerID := c.ComposerID
	ok = assert.NotEmpty(t, composerID, "composer id not provided")
	if !ok {
		return
	}

	composer, ok := GetComposer(composerID)
	if !ok {
		return assert.Fail(t, fmt.Sprintf("composer not found: %q", composerID))
	}

	images, ok := loadInputImages(t, c.InputImageNames...)
	if !ok {
		return ok
	}

	dc := gg.NewContext(c.ContextWidth(), c.ContextHeight())
	err := composer.Compose(dc, images...)
	ok = assert.NoError(t, err, "composer returned error")
	if !ok {
		return
	}

	expected, ok := c.ExpectedImage(t)
	if !ok {
		c.saveActualImage(t, dc)
		return
	}

	ok = assertImagesEqual(t, expected, dc.Image())
	if !ok {
		c.saveActualImage(t, dc)
	}

	return
}

var composerTests = []*ComposerTest{
	{
		ComposerID: "circles-pie",
		InputImageNames: []string{
			"b-martinez-744134.jpg",
			"j-crop-764891.jpg",
			"j-han-456323.jpg",
		},
	},
	{
		ComposerID: "tiles-perfect",
		InputImageNames: []string{
			"m-spiske-78531.jpg",
			"j-han-456323.jpg",
			"s-imbrock-487035.jpg",
			"s-erixon-753182.jpg",
		},
	},
	{
		ComposerID: "tiles-perfect",
		InputImageNames: []string{
			"s-erixon-753182.jpg",
			"j-han-456323.jpg",
			"m-spiske-78531.jpg",
			"b-martinez-744134.jpg",
			"s-imbrock-487035.jpg",
			"j-crop-764891.jpg",
		},
	},
	{
		ComposerID: "tiles-focused",
		InputImageNames: []string{
			"s-erixon-753182.jpg",
			"j-han-456323.jpg",
			"m-spiske-78531.jpg",
		},
	},
	{
		ComposerID: "tiles-focused",
		InputImageNames: []string{
			"s-erixon-753182.jpg",
			"j-han-456323.jpg",
			"m-spiske-78531.jpg",
			"b-martinez-744134.jpg",
			"s-imbrock-487035.jpg",
		},
	},
	{
		ComposerID: "tiles-focused",
		InputImageNames: []string{
			"s-erixon-753182.jpg",
			"j-han-456323.jpg",
			"m-spiske-78531.jpg",
			"b-martinez-744134.jpg",
			"s-imbrock-487035.jpg",
			"j-crop-764891.jpg",
		},
	},
	{
		ComposerID: "stripes-vertical",
		InputImageNames: []string{
			"j-crop-764891.jpg",
			"s-imbrock-487035.jpg",
			"s-erixon-753182.jpg",
		},
	},
}

func TestComposers(t *testing.T) {
	for _, c := range composerTests {
		t.Run(c.TestName(), func(t *testing.T) {
			c.Test(t)
		})
	}
}
