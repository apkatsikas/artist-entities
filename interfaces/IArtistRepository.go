package interfaces

import (
	"github.com/apkatsikas/artist-entities/models"
)

type IArtistRepository interface {
	Get(id uint) (*models.Artist, error)
	Create(artist *models.Artist) (*models.Artist, error)
	GetCount() (uint, error)
	GetByOffset(offset uint) (*models.Artist, error)
	Migrate() error
}
