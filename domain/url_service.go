package domain

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/domain/generator"
	"github.com/miladbarzideh/shortify/infra"
)

const maxRetries = 5

type Service struct {
	logger *logrus.Logger
	cfg    *infra.Config
	repo   *Repository
}

func NewService(logger *logrus.Logger, cfg *infra.Config, repo *Repository) *Service {
	return &Service{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

func (svc *Service) CreateShortURL(longURL string) (string, error) {
	shortCode := generator.GenerateShortURLCode(svc.cfg.Shortener.CodeLength)
	var url URL
	for i := 0; i < maxRetries; i++ {
		var err error
		url, err = svc.repo.Create(longURL, shortCode)
		if err == nil {
			break
		}

		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", err
		}

		svc.logger.Debugf("Failed to create short URL '%s'. Retrying...", longURL)
	}

	shortURL := svc.buildShortURL(url.ShortCode)
	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortURL":    shortURL,
	}).Debug("Create short URL")

	return shortURL, nil
}

func (svc *Service) GetLongURL(shortCode string) (string, error) {
	shortURL := svc.buildShortURL(shortCode)
	url, err := svc.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortCode":   shortURL,
	}).Debug("Get long URL")

	return url.LongURL, nil
}

func (svc *Service) buildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/api/v1/urls/%s", svc.cfg.Server.Address, shortCode)
}
