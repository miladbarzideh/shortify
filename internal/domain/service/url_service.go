package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	repository2 "github.com/miladbarzideh/shortify/internal/domain/repository"
	infra2 "github.com/miladbarzideh/shortify/internal/infra"
	"github.com/miladbarzideh/shortify/pkg/generator"
	"github.com/miladbarzideh/shortify/pkg/worker"
)

var (
	ErrURLNotFound = errors.New("url not found")
)

const maxRetries = 5

type URLService interface {
	CreateShortURL(ctx context.Context, url string) (string, error)
	GetLongURL(ctx context.Context, shortCode string) (string, error)
	BuildShortURL(shortCode string) string
	CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) error
}

type service struct {
	logger     *logrus.Logger
	cfg        *infra2.Config
	repo       repository2.URLRepository
	cacheRepo  repository2.URLCacheRepository
	gen        generator.Generator
	wp         worker.Pool
	cacheStats infra2.CacheStats
}

func NewService(logger *logrus.Logger,
	cfg *infra2.Config,
	repo repository2.URLRepository,
	cacheRepo repository2.URLCacheRepository,
	gen generator.Generator,
	wp worker.Pool,
	telemetry *infra2.TelemetryProvider,
) URLService {
	meter := telemetry.MeterProvider.Meter("urlService")
	return &service{
		logger:     logger,
		cfg:        cfg,
		repo:       repo,
		cacheRepo:  cacheRepo,
		gen:        gen,
		wp:         wp,
		cacheStats: infra2.NewCacheStats(meter),
	}
}

func (svc *service) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	shortCode := svc.gen.GenerateShortURLCode()
	if err := svc.wp.Submit(func() {
		if err := svc.CreateShortURLWithRetries(ctx, longURL, shortCode); err != nil {
			svc.logger.Error(err.Error())
		}
	}); err != nil {
		svc.logger.Error(err.Error())
		return "", err
	}

	shortURL := svc.BuildShortURL(shortCode)
	svc.logger.WithFields(logrus.Fields{
		"originalURL": longURL,
		"shortURL":    shortURL,
	}).Debug("Create short URL")

	return shortURL, nil
}

func (svc *service) CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) error {
	for i := 0; i < maxRetries; i++ {
		url, err := svc.repo.Create(ctx, longURL, shortCode)
		if err == nil {
			svc.logger.WithFields(logrus.Fields{
				"originalURL": url.LongURL,
				"shortCode":   url.ShortCode,
			}).Debug("short url created asynchronously")

			return nil
		}

		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return err
		}

		svc.logger.Debugf("failed to create short URL '%s'. Retrying...", longURL)
	}

	return fmt.Errorf("failed to create short URL after %d retries", maxRetries)
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