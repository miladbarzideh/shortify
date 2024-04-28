package service

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

const baseURL = "http://localhost:8080/api/v1"

var (
	// counter simulates a unique ID generator, similar to database IDs
	counter atomic.Uint64
)

type Service struct {
	logger    *logrus.Logger
	shortURLs map[string]string
	longURLs  map[string]string
}

func NewService(logger *logrus.Logger) *Service {
	return &Service{
		logger:    logger,
		shortURLs: make(map[string]string),
		longURLs:  make(map[string]string),
	}
}

func (svc *Service) CreateShortURL(longURL string) (string, error) {
	if shortUrl, exists := svc.longURLs[longURL]; exists {
		return shortUrl, nil
	}

	// To concurrency issue database integration can provide a solution
	id := counter.Add(1)
	shortCode := base64.URLEncoding.EncodeToString([]byte(strconv.Itoa(int(id))))
	svc.shortURLs[shortCode] = longURL
	svc.longURLs[longURL] = shortCode
	shortURL := buildShortURL(shortCode)
	svc.logger.WithFields(logrus.Fields{
		"originalURL": longURL,
		"shortURL":    shortURL,
	}).Debug("Create short URL")

	return shortURL, nil
}

func (svc *Service) GetLongURL(shortCode string) (string, error) {
	shortURL := buildShortURL(shortCode)
	if longURL, exists := svc.shortURLs[shortCode]; exists {
		svc.logger.WithFields(logrus.Fields{
			"originalURL": longURL,
			"shortCode":   shortURL,
		}).Debug("Get long URL")
		return longURL, nil
	}

	return "", fmt.Errorf("short url not found for %s", shortURL)
}

func buildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortCode)
}
