package mosaic

import (
	"github.com/fogleman/gg"
	"image"
)

// A Composer creates image compositions
type Composer interface {
	// Compose draws the images to the drawing context.
	Compose(dc *gg.Context, images ...image.Image) error
}

// A ComposerFunc is a Composer which itself is a function.
type ComposerFunc func(dc *gg.Context, images ...image.Image) error

func (f ComposerFunc) Compose(dc *gg.Context, images ...image.Image) error {
	return f(dc, images...)
}

// TODO:
//  - recommended amount of images,
//  - required amount of images (function check and human repr)
//  - description?

// ComposerInfo is a Composer with additional information.
type ComposerInfo struct {
	Composer

	Id   string
	Name string
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
