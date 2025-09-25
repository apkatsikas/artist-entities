package interfaces

import (
	"github.com/apkatsikas/artist-entities/models"
)

type IArtistService interface {
	Get(id uint) (*models.Artist, error)
	Create(artistName string) (*models.Artist, error)
	GetRandom() (*models.Artist, error)
}
