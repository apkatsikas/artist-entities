package repositories

import (
	"errors"

	ce "github.com/apkatsikas/artist-entities/customerrors"
	"github.com/apkatsikas/artist-entities/interfaces"
	"github.com/apkatsikas/artist-entities/models"
	"gorm.io/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ArtistRepository struct {
	IDB interfaces.IDbHandler
}

func (ar *ArtistRepository) GetCount() (uint, error) {
	gormConn := ar.IDB.Connection()

	var count int64

	result := gormConn.Model(models.Artist{}).Count(&count)

	if result.Error != nil {
		return uint(0), result.Error
	}
	return uint(count), nil
}

func (ar *ArtistRepository) GetByOffset(offset uint) (*models.Artist, error) {
	gormConn := ar.IDB.Connection()

	var artist = models.Artist{}

	result := gormConn.Limit(1).Offset(int(offset)).Find(&artist)

	if result.Error != nil {
		return nil, result.Error
	}

	return &artist, nil
}

func (ar *ArtistRepository) Get(id uint) (*models.Artist, error) {
	var artist = models.Artist{}
	artist.ID = id
	result := ar.IDB.Connection().First(&artist)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ce.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &artist, nil
}

func (ar *ArtistRepository) Create(artist *models.Artist) (*models.Artist, error) {
	// If record can't be found, insert it
	result := ar.IDB.Connection().Where(
		models.Artist{Name: artist.Name}).FirstOrCreate(&artist)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ce.ErrRecordExists
	}
	return artist, nil
}

func (ar *ArtistRepository) Migrate() error {
	// Create table if needed
	err := ar.IDB.Connection().AutoMigrate(&models.Artist{})
	if err != nil {
		return err
	}

	return nil
}
