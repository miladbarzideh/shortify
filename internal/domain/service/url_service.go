package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/domain/repository"
	"github.com/miladbarzideh/shortify/internal/infra"
	"github.com/miladbarzideh/shortify/pkg/generator"
)

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
)

const maxRetries = 5

type URLService interface {
	CreateShortURL(ctx context.Context, url string) (string, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	BuildShortURL(shortCode string) string
	CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) (*model.URL, error)
}

type service struct {
	logger     *logrus.Logger
	cfg        *infra.Config
	repo       repository.URLRepository
	cacheRepo  repository.URLCacheRepository
	gen        generator.Generator
	cacheStats infra.CacheStats
}

func NewService(logger *logrus.Logger,
	cfg *infra.Config,
	repo repository.URLRepository,
	cacheRepo repository.URLCacheRepository,
	gen generator.Generator,
	telemetry *infra.TelemetryProvider,
) URLService {
	meter := telemetry.MeterProvider.Meter("urlService")
	return &service{
		logger:     logger,
		cfg:        cfg,
		repo:       repo,
		cacheRepo:  cacheRepo,
		gen:        gen,
		cacheStats: infra.NewCacheStats(meter),
	}
}

func (svc *service) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	shortCode := svc.gen.GenerateShortURLCode()
	url, err := svc.CreateShortURLWithRetries(ctx, longURL, shortCode)
	if err != nil {
		return "", err
	}

	shortURL := svc.BuildShortURL(url.ShortCode)
	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortURL":    shortURL,
	}).Debug("Create short URL")

	return shortURL, nil
}

func (svc *service) CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) (*model.URL, error) {
	url := &model.URL{ShortCode: shortCode, LongURL: longURL}
	for i := 0; i < maxRetries; i++ {
		err := svc.repo.Create(ctx, url)
		if err == nil {
			return url, nil
		}

		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, err
		}

		svc.logger.Debugf("failed to create short URL '%s'. Retrying...", longURL)
	}

	return nil, fmt.Errorf("failed to create short URL after %d retries %w", maxRetries, ErrMaxRetriesExceeded)
}

func (svc *service) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	if url, err := svc.cacheRepo.Get(ctx, shortCode); err == nil {
		svc.cacheStats.Hits.Inc(ctx)
		return url.LongURL, nil
	}

	url, err := svc.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrURLNotFound
		}

		return "", err
	}

	if err = svc.cacheRepo.Set(ctx, url); err != nil {
		svc.logger.Errorf("failed to cache short URL '%s'. Error: %v", shortCode, err)
	}

	svc.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortURL":    svc.BuildShortURL(shortCode),
	}).Debug("read URL from database")
	svc.cacheStats.Hits.Inc(ctx)

	return url.LongURL, nil
}

func (svc *service) BuildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/api/v1/urls/%s", svc.cfg.Server.Address, shortCode)
}
