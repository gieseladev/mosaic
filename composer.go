package mosaic

import (
	"github.com/fogleman/gg"
	"image"
	"sort"
)

// A Composer creates image compositions
type Composer interface {
	// Compose draws the images to the drawing context.
	Compose(dc *gg.Context, images ...image.Image) error
}

// A ComposerFunc is a Composer which itself is a function.
type ComposerFunc func(dc *gg.Context, images ...image.Image) error

// Compose calls the underlying function with the given arguments.
func (f ComposerFunc) Compose(dc *gg.Context, images ...image.Image) error {
	return f(dc, images...)
}

// ComposerInfo is a Composer with additional information.
type ComposerInfo struct {
	Composer

	Id          string
	Name        string
	Description string

	ImageCountHuman string
	CheckImageCount func(count int) bool

	RecommendedImageCounts []int
}

// RecommendImageCount recommends a suitable amount of images to use
// which is guaranteed to be less or equal to the amount provided.
func (ci ComposerInfo) RecommendImageCount(imageCount int) int {
	var currentBest int
	for _, c := range ci.RecommendedImageCounts {
		if c > currentBest && c <= imageCount {
			currentBest = c
		}
	}

	// found a recommended count
	if currentBest > 0 {
		return currentBest
	}

	// no check provided, assume all values are possible
	if ci.CheckImageCount == nil {
		return imageCount
	}

	// check all counts and take the first that works
	for i := imageCount; i > 0; i-- {
		if ci.CheckImageCount(i) {
			return i
		}
	}

	return 0
}

var registeredComposers []ComposerInfo

// RegisterComposer registers the given
func RegisterComposer(comps ...ComposerInfo) error {
	for _, c := range comps {
		registeredComposers = append(registeredComposers, c)
	}

	return nil
}

// GetComposer returns the composer with the given id.
func GetComposer(id string) (ComposerInfo, bool) {
	for _, composer := range registeredComposers {
		if composer.Id == id {
			return composer, true
		}
	}

	return ComposerInfo{}, false
}

// GetComposers returns a slice containing all composers.
func GetComposers() []ComposerInfo {
	return registeredComposers
}

// RecommendComposers returns a slice of composers which are suitable
// for the given image count.
func RecommendComposers(count int) []ComposerInfo {
	type ComposerComp struct {
		C             ComposerInfo
		RecImageCount int
	}

	composerComparisons := make([]ComposerComp, 0)

	for _, composer := range registeredComposers {
		recommended := composer.RecommendImageCount(count)
		if recommended != 0 {
			composerComparisons = append(composerComparisons, ComposerComp{
				composer,
				recommended,
			})
		}

	}

	sort.Slice(composerComparisons, func(i, j int) bool {
		return composerComparisons[i].RecImageCount < composerComparisons[j].RecImageCount
	})

	composers := make([]ComposerInfo, len(composerComparisons))
	for i, cc := range composerComparisons {
		composers[i] = cc.C
	}

	return composers
}
