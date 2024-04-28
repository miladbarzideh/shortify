package service

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/sirupsen/logrus"
)

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
	shortURL := base64.URLEncoding.EncodeToString([]byte(strconv.Itoa(int(id))))
	svc.shortURLs[shortURL] = longURL
	svc.longURLs[longURL] = shortURL

	return shortURL, nil
}

func (svc *Service) GetLongURL(shortURL string) (string, error) {
	if longURL, exists := svc.shortURLs[shortURL]; exists {
		return longURL, nil
	}

	return "", fmt.Errorf("short url not found for %s", shortURL)
}
