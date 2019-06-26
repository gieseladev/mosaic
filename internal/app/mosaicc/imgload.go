package mosaicc

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func loadImageFromURL(u string) (image.Image, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, nil
	}

	img, _, err := image.Decode(resp.Body)
	_ = resp.Body.Close()
	return img, err
}

// LoadImage loads an image from the given location.
// The location can be either a url, or a filepath pointing
// to an image.
func LoadImage(location string) (image.Image, error) {
	_, err := url.ParseRequestURI(location)
	if err == nil {
		return loadImageFromURL(location)
	} else {
		return gg.LoadImage(location)
	}
}

// LoadImages loads the given images in parallel.
func LoadImages(locations []string) ([]image.Image, error) {
	type LoadResult struct {
		Index int
		Image image.Image
		Err   error
	}

	resultChan := make(chan LoadResult)

	lastIndex := len(locations) - 1
	for i, location := range locations {
		go func(i int, location string) {
			img, err := LoadImage(location)
			if err != nil {
				err = fmt.Errorf("couldn't load image %q: %v", location, err)
			}

			resultChan <- LoadResult{
				Index: i,
				Image: img,
				Err:   err,
			}

			if i == lastIndex {
				close(resultChan)
			}
		}(i, location)
	}

	errs := make([]error, 0)
	results := make([]LoadResult, 0, len(locations))
	for result := range resultChan {
		if result.Err == nil {
			results = append(results, result)
		} else {
			errs = append(errs, result.Err)
		}
	}

	// preserve original order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Index < results[j].Index
	})

	images := make([]image.Image, len(results))
	for i, result := range results {
		images[i] = result.Image
	}

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
