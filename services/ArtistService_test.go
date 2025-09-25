package services

import (
	"errors"
	"testing"

	ce "github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/interfaces/mocks"
	"github.com/apkatsikas/artist-entities/models"

	"github.com/stretchr/testify/assert"
)

const (
	weirdError = "weird error"
)

type artistServiceTestMocks struct {
	*mocks.IArtistRules
	*mocks.IArtistRepository
}

func artistServiceReqMocks(t *testing.T) artistServiceTestMocks {
	return artistServiceTestMocks{
		IArtistRules:      mocks.NewIArtistRules(t),
		IArtistRepository: mocks.NewIArtistRepository(t),
	}
}

func injectedArtistService(mocks artistServiceTestMocks) ArtistService {
	return ArtistService{
		ArtistRepository: mocks.IArtistRepository,
		Rules:            mocks.IArtistRules,
	}
}

func TestGetArtist(t *testing.T) {
	// Artist data
	artistName := "Lou Reed"
	artistID := uint(1)
	artist := models.Artist{}
	artist.ID = artistID
	artist.Name = artistName

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRepository.EXPECT().Get(artistID).Return(&artist, nil)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Get artist
	artistResult, err := artistService.Get(artistID)

	// Check artist result
	assert.Equal(t, &artist, artistResult)

	// Check that there is no error
	assert.Nil(t, err)
}

func TestGetArtistNoRecord(t *testing.T) {
	// Artist data
	artistID := uint(2)

	// Expected error
	expectedError := ce.ErrRecordNotFound

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRepository.EXPECT().Get(artistID).Return(nil, expectedError)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Get artist
	artistResult, err := artistService.Get(artistID)

	// Check that we got no artist
	assert.Nil(t, artistResult)

	// Check error
	assert.True(t, errors.Is(err, expectedError))
}

func TestCreateArtist(t *testing.T) {
	// Artist data
	artistName := "Black Sabbath"
	artistID := uint(420)
	artist := models.Artist{}
	artist.Name = artistName
	artist.ID = artistID

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRules.EXPECT().CleanArtistName(artistName).Return(artistName, nil)
	mocks.IArtistRepository.EXPECT().Create(artistName).Return(&artist, nil)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Create artist
	artistResult, err := artistService.Create(artistName)

	// Check artist result
	assert.Equal(t, &artist, artistResult)

	// Check error
	assert.Nil(t, err)
}

func TestCreateArtistExists(t *testing.T) {
	// Artist data
	artistName := "Lou Reed"

	// Expected error
	expectedError := ce.ErrRecordExists

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRules.EXPECT().CleanArtistName(artistName).Return(artistName, nil)
	mocks.IArtistRepository.EXPECT().Create(artistName).Return(nil, expectedError)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Create artist
	artistResult, err := artistService.Create(artistName)

	// Check that we got no artist
	assert.Nil(t, artistResult)

	// Check error
	assert.True(t, errors.Is(err, expectedError))
}

func TestCreateArtistRulesFail(t *testing.T) {
	// Artist data
	artistName := "Black,Sabbath"

	// Expected error
	expectedError := ce.ErrDataInvalid

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRules.EXPECT().CleanArtistName(artistName).Return("", expectedError)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Create artist
	artistResult, err := artistService.Create(artistName)

	// Check that we got no artist
	assert.Nil(t, artistResult)

	// Check error
	assert.True(t, errors.Is(err, expectedError))
}

func TestGetRandomArtist(t *testing.T) {
	// Setup data
	count := uint(666)
	offsetValue := uint(187)

	// Artist data
	artistName := "Black Sabbath"
	artistID := offsetValue
	artist := models.Artist{}
	artist.Name = artistName
	artist.ID = artistID

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRepository.EXPECT().GetCount().Return(count, nil)
	mocks.IArtistRules.EXPECT().RandomOffset(count).Return(offsetValue)
	mocks.IArtistRepository.EXPECT().GetByOffset(offsetValue).Return(&artist, nil)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Get artist
	artistResult, err := artistService.GetRandom()

	// Check artist
	assert.Equal(t, &artist, artistResult)

	// Check error
	assert.Nil(t, err)
}

func TestGetRandomArtistCountError(t *testing.T) {
	// Expected error
	expectedError := errors.New(weirdError)

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRepository.EXPECT().GetCount().Return(uint(0), expectedError)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Get artist
	artistResult, err := artistService.GetRandom()

	// Check that we got no artist
	assert.Nil(t, artistResult)
	// Check error
	assert.True(t, errors.Is(err, expectedError))
}

func TestGetRandomArtistOffsetError(t *testing.T) {
	// Setup data
	count := uint(666)
	offsetValue := uint(187)

	// Expected error
	expectedError := errors.New(weirdError)

	// Setup mocks
	mocks := artistServiceReqMocks(t)
	mocks.IArtistRepository.EXPECT().GetCount().Return(count, nil)
	mocks.IArtistRules.EXPECT().RandomOffset(count).Return(offsetValue)
	mocks.IArtistRepository.EXPECT().GetByOffset(offsetValue).Return(nil, expectedError)

	// Inject service
	artistService := injectedArtistService(mocks)

	// Get artist
	artistResult, err := artistService.GetRandom()

	// Check that we got no artist
	assert.Nil(t, artistResult)
	// Check error
	assert.True(t, errors.Is(err, expectedError))
}
