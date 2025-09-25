package interfaces

import (
	"github.com/apkatsikas/artist-entities/models"
)

type IUserRepository interface {
	Get(name string) (*models.User, error)
	Create(name string, password string) (*models.User, error)
}
