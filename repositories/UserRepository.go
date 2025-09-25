package repositories

import (
	"errors"

	ce "github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/models"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UserRepository struct {
	IDB interfaces.IDbHandler
}

func (ur *UserRepository) Get(name string) (*models.User, error) {
	var user = models.User{}
	user.Name = name
	result := ur.IDB.Connection().First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ce.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) Create(name string, password string) (*models.User, error) {
	var user models.User

	result := ur.IDB.Connection().Where(
		models.User{Name: name, Password: password}).FirstOrCreate(&user)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ce.ErrRecordExists
	}
	return &user, nil
}

func (ur *UserRepository) Migrate() error {
	err := ur.IDB.Connection().AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
