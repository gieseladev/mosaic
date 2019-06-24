package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindBalancedFactors(t *testing.T) {
	a, b := FindBalancedFactors(9)
	assert.Equal(t, 3+3, a+b)

	a, b = FindBalancedFactors(11)
	assert.Equal(t, 1+11, a+b)

	a, b = FindBalancedFactors(15)
	assert.Equal(t, 3+5, a+b)

	a, b = FindBalancedFactors(6)
	assert.Equal(t, 2+3, a+b)
}
