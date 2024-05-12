package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/model"
)

type URLRepository interface {
	Create(longURL string, shortCode string) (model.URL, error)
	FindByShortCode(shortCode string) (model.URL, error)
}

type Repository struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewRepository(logger *logrus.Logger, db *gorm.DB) URLRepository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}

func (r Repository) Create(longURL string, shortCode string) (model.URL, error) {
	url := model.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	result := r.db.Create(&url)

	return url, result.Error
}

func (r Repository) FindByShortCode(shortCode string) (model.URL, error) {
	var url model.URL
	result := r.db.Where("short_code = ?", shortCode).First(&url)

	return url, result.Error
}
