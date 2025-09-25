package services

import (
	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/models"
)

type ArtistService struct {
	ArtistRepository interfaces.IArtistRepository
	Rules            interfaces.IArtistRules
}

func (as *ArtistService) Get(id uint) (*models.Artist, error) {
	// Get artist
	artist, err := as.ArtistRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (as *ArtistService) GetRandom() (*models.Artist, error) {
	count, err := as.ArtistRepository.GetCount()

	if err != nil {
		return nil, err
	}

	offset := as.Rules.RandomOffset(count)

	artist, err := as.ArtistRepository.GetByOffset(offset)

	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (as *ArtistService) Create(artistName string) (*models.Artist, error) {
	// Clean
	artistName, err := as.Rules.CleanArtistName(artistName)
	if err != nil {
		return nil, err
	}

	// Write to repository
	artist, err := as.ArtistRepository.Create(artistName)
	if err != nil {
		return nil, err
	}

	return artist, nil
}
