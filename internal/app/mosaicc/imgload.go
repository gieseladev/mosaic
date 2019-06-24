package mosaicc

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"strings"
	"sync"
)

func LoadImage(location string) (image.Image, error) {
	// TODO handle urls
	return gg.LoadImage(location)
}

func LoadImages(locations []string) ([]image.Image, error) {
	var mut sync.Mutex
	var wg sync.WaitGroup

	wg.Add(len(locations))

	var images []image.Image
	var errs []error
	for _, location := range locations {
		go func(location string) {
			img, err := LoadImage(location)

			mut.Lock()
			defer mut.Unlock()

			if err == nil {
				images = append(images, img)
			} else {
				errs = append(errs, fmt.Errorf("couldn't load image %q: %v", location, err))
			}

			wg.Done()
		}(location)
	}

	wg.Wait()

	var combinedErr error
	if len(errs) > 0 {
		errStrings := make([]string, len(errs))
		for i, err := range errs {
			errStrings[i] = err.Error()
		}

		combinedErr = fmt.Errorf("some images couldn't be loaded:\n%s", strings.Join(errStrings, "\n"))
	}

	return images, combinedErr
}
