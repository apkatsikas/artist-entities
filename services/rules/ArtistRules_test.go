package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomOffset(t *testing.T) {
	count := uint(666)
	rules := ArtistRules{}
	result := rules.RandomOffset(count)

	withinBounds := result >= 1 && result <= count
	assert.True(t, withinBounds)
}

func TestRandomCount1(t *testing.T) {
	count := uint(1)
	rules := ArtistRules{}

	result := rules.RandomOffset(count)

	assert.Equal(t, count, result)
}

func TestRandomOffsetPanics(t *testing.T) {
	count := uint(0)
	rules := ArtistRules{}

	assert.Panics(t, func() { rules.RandomOffset(count) })
}
