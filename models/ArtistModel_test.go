package models_test

import (
	"errors"
	"testing"

	"github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/models"
	"github.com/stretchr/testify/assert"
)

const (
	cleanSabbath = "black sabbath"
	cleanStones  = "rolling stones"
	sly          = "sly and the family stone"
)

func TestArtistCleanName(t *testing.T) {
	var testData = []struct {
		test     string
		name     string
		expected string
	}{
		{test: "pass through", name: "black sabbath", expected: cleanSabbath},
		{test: "leading whitespace", name: " black sabbath", expected: cleanSabbath},
		{test: "trailing whitespace", name: "black sabbath ", expected: cleanSabbath},
		{test: "leading AND trailing whitespace", name: " black sabbath ", expected: cleanSabbath},
		{test: "caps", name: "Black Sabbath", expected: cleanSabbath},
		{test: "leading the", name: "the rolling stones", expected: cleanStones},
		{test: "fake the", name: "them", expected: "them"},
		{test: "the the", name: "the the", expected: "the"},
		{test: "sly", name: sly, expected: sly},
		{test: "one single quote", name: "howlin' wolf", expected: "howlin wolf"},
		{test: "period", name: "n.w.a", expected: "nwa"},
		{test: "exclamation", name: "neu!", expected: "neu"},
		{test: "the works", name: " the Rolling' Stones ", expected: cleanStones},
		{test: "cutting the makes them pass",
			name:     "the insanelylongnamecanubelievethistrulyincrediblewhatkindoflamebandwouldhave",
			expected: "insanelylongnamecanubelievethistrulyincrediblewhatkindoflamebandwouldhave"},
	}
	for _, tt := range testData {
		t.Run(tt.test, func(t *testing.T) {
			result, err := models.ValidatedArtist(tt.name)

			assert.Equal(t, tt.expected, result.Name)
			assert.Nil(t, err)
		})
	}
}

func TestArtistCleanLongName(t *testing.T) {
	artistName := "insanelylongnamecanubelievethistrulyincrediblewhatkindoflamebandwouldhavethis"
	result, err := models.ValidatedArtist(artistName)

	assert.True(t, errors.Is(err, customerrors.ErrDataTooLong))
	assert.Empty(t, result)
}

func TestArtistInvalidData(t *testing.T) {
	var testData = []struct {
		test string
		name string
	}{
		{test: "commas", name: "how,happen"},
		{test: "2 single quotes", name: "'blargh'"},
		{test: "double quotes", name: "\""},
		{test: "the works", name: "wtf,'\"isthis"},
		{test: "chk chk chk", name: "!!!"},
	}
	for _, tt := range testData {
		t.Run(tt.test, func(t *testing.T) {
			result, err := models.ValidatedArtist(tt.name)

			assert.True(t, errors.Is(err, customerrors.ErrDataInvalid))
			assert.Empty(t, result)
		})
	}
}
