package domain

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	logger *logrus.Logger
	db     *gorm.DB
}

func NewRepository(logger *logrus.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}

func (r Repository) Create(longURL string, shortCode string) (URL, error) {
	url := URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	result := r.db.Create(&url)

	return url, result.Error
}

func (r Repository) FindByShortCode(shortCode string) (URL, error) {
	var url URL
	result := r.db.Where("short_code = ?", shortCode).First(&url)

	return url, result.Error
}
