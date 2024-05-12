package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/infra"
	"github.com/miladbarzideh/shortify/internal/model"
	"github.com/miladbarzideh/shortify/internal/repository"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

var (
	ErrURLNotFound = errors.New("url not found")
)

const maxRetries = 5

type URLService interface {
	CreateShortURL(url string) (string, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	BuildShortURL(shortCode string) string
}

type service struct {
	logger    *logrus.Logger
	cfg       *infra.Config
	repo      repository.URLRepository
	cacheRepo repository.URLCacheRepository
	gen       generator.Generator
}

func NewService(
	logger *logrus.Logger,
	cfg *infra.Config,
	repo repository.URLRepository,
	cacheRepo repository.URLCacheRepository,
	gen generator.Generator,
) URLService {
	return &service{
		logger:    logger,
		cfg:       cfg,
		repo:      repo,
		cacheRepo: cacheRepo,
		gen:       gen,
	}
}

func (svc *service) CreateShortURL(longURL string) (string, error) {
	shortCode := svc.gen.GenerateShortURLCode(svc.cfg.Shortener.CodeLength)
	var url model.URL
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

	shortURL := svc.BuildShortURL(url.ShortCode)
	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortURL":    shortURL,
	}).Debug("Create short URL")

	return shortURL, nil
}

func (svc *service) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	if url, err := svc.cacheRepo.Get(ctx, shortCode); err == nil {
		return url.LongURL, nil
	}

	url, err := svc.repo.FindByShortCode(shortCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	if err = svc.cacheRepo.Set(ctx, url); err != nil {
		svc.logger.Errorf("Failed to cache short URL '%s'. Error: %v", shortCode, err)
	}

	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortURL":    svc.BuildShortURL(shortCode),
	}).Debug("Read URL from database")

	return url.LongURL, nil
}

func (svc *service) BuildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/api/v1/urls/%s", svc.cfg.Server.Address, shortCode)
}
