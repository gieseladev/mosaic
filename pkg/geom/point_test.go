package geom

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPoint_Rotate(t *testing.T) {
	p := Pt(1, 0).Rotate(QuarterPi)
	assert.InDelta(t, 1/math.Sqrt2, p.X, .0000001)
	assert.InDelta(t, 1/math.Sqrt2, p.Y, .0000001)
}
