package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/infra"
)

const (
	cachePrefix = "short-url"
	cacheTTL    = 24 * time.Hour
)

type CacheRepository struct {
	logger *logrus.Logger
	cache  *redis.Client
	tracer trace.Tracer
}

func NewCacheRepository(logger *logrus.Logger, redis *redis.Client, telemetry *infra.TelemetryProvider) *CacheRepository {
	tracer := telemetry.TraceProvider.Tracer("urlCacheRepo")
	return &CacheRepository{
		logger: logger,
		cache:  redis,
		tracer: tracer,
	}
}

func (cr *CacheRepository) Set(ctx context.Context, url *model.URL) error {
	_, span := cr.tracer.Start(ctx, "urlCacheRepo.set")
	defer span.End()
	value, err := json.Marshal(url)
	if err != nil {
		return err
	}

	err = cr.cache.Set(ctx, cr.buildKeyWithPrefix(url.ShortCode), value, cacheTTL).Err()
	if err != nil {
		return err
	}

	cr.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortCode":   url.ShortCode,
	}).Debug("Write URL to cache")

	return nil
}

func (cr *CacheRepository) Get(ctx context.Context, shortCode string) (*model.URL, error) {
	_, span := cr.tracer.Start(ctx, "urlCacheRepo.get")
	defer span.End()
	var url model.URL
	result, err := cr.cache.Get(ctx, cr.buildKeyWithPrefix(shortCode)).Result()
	if err != nil {
		cr.logger.Error(err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &url); err != nil {
		cr.logger.Error(err)
		return nil, err
	}

	cr.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortCode":   shortCode,
	}).Debug("Read URL from cache")

	return &url, nil
}

func (cr *CacheRepository) buildKeyWithPrefix(url string) string {
	return fmt.Sprintf("%s:%s", cachePrefix, url)
}
