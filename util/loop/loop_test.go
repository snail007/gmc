package gloop

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFor(t *testing.T) {
	t.Parallel()
	i := 0
	For(10, func(idx int) {
		i += idx
	})
	assert.Equal(t, 45, i)
}

func TestFor1(t *testing.T) {
	t.Parallel()
	i := 0
	j := 0
	k := 0
	ForBy(9, -1, func(idx, value int) {
		i += idx
		j += value
		k += value - idx
	})
	assert.Equal(t, 45, i)
	assert.Equal(t, 45, j)
	assert.Equal(t, 0, k)
}

func TestDoWhile(t *testing.T) {
	t.Parallel()
	i := 0
	Do(func(idx int) {
		i += idx
	}).While(func(idx int) bool {
		return idx < 10
	})
	assert.Equal(t, 45, i)
}

func TestDoUntil(t *testing.T) {
	t.Parallel()
	i := 0
	Do(func(idx int) {
		i += idx
	}).Until(func(idx int) bool {
		return idx > 9
	})
	assert.Equal(t, 45, i)
}
