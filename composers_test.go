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

func loadInputImages(t testing.TB, names ...string) ([]image.Image, bool) {
	images := make([]image.Image, len(names))

	for i := 0; i < len(names); i++ {
		img, err := loadImage(fmt.Sprintf("test/data/input/%s", names[i]))
		ok := assert.NoError(t, err, fmt.Sprintf("couldn't load input images: %v", names))
		if !ok {
			return images, ok
		}

		images[i] = img
	}

	return images, true
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

func (c *ComposerTest) GetComposer(t testing.TB) (composer Composer, ok bool) {
	composerID := c.ComposerID
	ok = assert.NotEmpty(t, composerID, "composer id not provided")
	if !ok {
		return
	}

	composer, ok = GetComposer(composerID)
	if !ok {
		ok = assert.Fail(t, fmt.Sprintf("composer not found: %q", composerID))
	}

	return
}

func (c *ComposerTest) saveActualImage(t *testing.T, dc *gg.Context) bool {
	err := dc.SavePNG(fmt.Sprintf("%s/%s-actual.png",
		"test/data/output",
		c.TestName(),
	))

	return assert.NoError(t, err, "couldn't save actual output")
}

func (c *ComposerTest) Test(t *testing.T) (ok bool) {
	composer, ok := c.GetComposer(t)
	if !ok {
		return
	}

	images, ok := loadInputImages(t, c.InputImageNames...)
	if !ok {
		return
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

func (c *ComposerTest) Benchmark(b *testing.B) {
	composer, ok := c.GetComposer(b)
	if !ok {
		return
	}

	images, ok := loadInputImages(b, c.InputImageNames...)
	if !ok {
		return
	}

	dc := gg.NewContext(c.ContextWidth(), c.ContextHeight())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dc.Clear()
		_ = composer.Compose(dc, images...)
	}
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
		ComposerID: "tiles-diamond",
		InputImageNames: []string{
			"b-martinez-744134.jpg",
			"i-palacio-Y20JJ_ddy9M.jpg",
			"j-crop-764891.jpg",
			"j-han-456323.jpg",
			"j-pereira-fSGsKbICefw.jpg",
			"j-wejxKZ-9IZg.jpg",
			"m-spiske-78531.jpg",
			"m-wingen-PDX_a_82obo.jpg",
			"n-perea-W8BRzoUTHNA.jpg",
			"p-wooten-FMiczIq8orU.jpg",
			"s-erixon-753182.jpg",
			"s-imbrock-487035.jpg",
			"t-mikuckis-hbnH0ILjUZE.jpg",
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
		c := c

		t.Run(c.TestName(), func(t *testing.T) {
			t.Parallel()
			c.Test(t)
		})
	}
}

func BenchmarkComposers(b *testing.B) {
	for _, c := range composerTests {
		b.Run(c.TestName(), func(b *testing.B) {
			c.Benchmark(b)
		})
	}
}
