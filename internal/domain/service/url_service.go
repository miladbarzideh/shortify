package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/infra"
)

var (
	ErrURLNotFound        = errors.New("url not found")
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
)

const maxRetries = 5

type URLRepository interface {
	Create(ctx context.Context, url *model.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
}

type URLCacheRepository interface {
	Set(ctx context.Context, url *model.URL) error
	Get(ctx context.Context, shortCode string) (*model.URL, error)
	BuildKeyWithPrefix(url string) string
}

type Generator interface {
	GenerateShortURLCode() string
}

type Service struct {
	logger     *logrus.Logger
	cfg        *infra.Config
	repo       URLRepository
	cacheRepo  URLCacheRepository
	gen        Generator
	cacheStats infra.CacheStats
}

func NewService(logger *logrus.Logger,
	cfg *infra.Config,
	repo URLRepository,
	cacheRepo URLCacheRepository,
	gen Generator,
	telemetry *infra.TelemetryProvider,
) *Service {
	meter := telemetry.MeterProvider.Meter("urlService")
	return &Service{
		logger:     logger,
		cfg:        cfg,
		repo:       repo,
		cacheRepo:  cacheRepo,
		gen:        gen,
		cacheStats: infra.NewCacheStats(meter),
	}
}

func (svc *Service) CreateShortURL(ctx context.Context, longURL string) (string, error) {
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

func (svc *Service) CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) (*model.URL, error) {
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

func (svc *Service) GetLongURL(ctx context.Context, shortCode string) (string, error) {
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

func (svc *Service) BuildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/api/v1/urls/%s", svc.cfg.Server.Address, shortCode)
}
