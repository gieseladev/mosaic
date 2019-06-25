package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolygon_Center(t *testing.T) {
	assert.Equal(t, Pt(1, 1), Poly(
		Pt(0, 2), Pt(2, 2),
		Pt(0, 0), Pt(2, 0),
	).Center())

	assert.Equal(t, Pt(1.5, 1.5), Poly(
		Pt(1, 2), Pt(2, 2),
		Pt(1, 1), Pt(2, 1),
	).Center())
}

func TestPolygon_Translate(t *testing.T) {
	assert.Equal(t, Poly(
		Pt(.5, 2.5), Pt(2.5, 2.5),
		Pt(.5, .5), Pt(2.5, .5),
	), Poly(
		Pt(1, 2), Pt(2, 2),
		Pt(1, 1), Pt(2, 1),
	).ScaleAroundCenter(2))
}
