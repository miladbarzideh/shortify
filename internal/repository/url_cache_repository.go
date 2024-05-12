package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/miladbarzideh/shortify/internal/model"
)

const (
	cachePrefix = "short-url"
	cacheTTL    = 24 * time.Hour
)

type URLCacheRepository interface {
	Set(ctx context.Context, url model.URL) error
	Get(ctx context.Context, shortCode string) (model.URL, error)
	BuildKeyWithPrefix(url string) string
}

type cacheRepository struct {
	logger *logrus.Logger
	cache  *redis.Client
}

func NewCacheRepository(logger *logrus.Logger, redis *redis.Client) URLCacheRepository {
	return &cacheRepository{
		logger: logger,
		cache:  redis,
	}
}

func (cr *cacheRepository) Set(ctx context.Context, url model.URL) error {
	value, err := json.Marshal(url)
	if err != nil {
		return err
	}

	err = cr.cache.Set(ctx, cr.BuildKeyWithPrefix(url.ShortCode), value, cacheTTL).Err()
	if err != nil {
		return err
	}

	cr.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortCode":   url.ShortCode,
	}).Debug("Write URL to cache")

	return nil
}

func (cr *cacheRepository) Get(ctx context.Context, shortCode string) (model.URL, error) {
	var url model.URL
	result, err := cr.cache.Get(ctx, cr.BuildKeyWithPrefix(shortCode)).Result()
	if err != nil {
		cr.logger.Error(err)
		return url, err
	}

	if err = json.Unmarshal([]byte(result), &url); err != nil {
		cr.logger.Error(err)
		return url, err
	}

	cr.logger.WithFields(logrus.Fields{
		"originalURL": url.LongURL,
		"shortCode":   shortCode,
	}).Debug("Read URL from cache")

	return url, nil
}

func (cr *cacheRepository) BuildKeyWithPrefix(url string) string {
	return fmt.Sprintf("%s:%s", cachePrefix, url)
}