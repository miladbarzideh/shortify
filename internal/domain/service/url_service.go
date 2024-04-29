package service

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

const baseURL = "http://localhost:8080/api/v1/urls"

var (
	// counter simulates a unique ID generator, similar to database IDs
	counter uint
	mutex   sync.Mutex
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
	mutex.Lock()
	defer mutex.Unlock()
	if shortCode, exists := svc.longURLs[longURL]; exists {
		return buildShortURL(shortCode), nil
	}

	counter++
	shortCode := Base62EncodeWithObfuscatedID(uint64(counter))
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

	return "", fmt.Errorf("long url not found for %s", shortURL)
}

func buildShortURL(shortCode string) string {
	return fmt.Sprintf("%s/%s", baseURL, shortCode)
}
