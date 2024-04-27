package service

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
}

func NewService(logger *logrus.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (svc *Service) CreateShortURL(longURL string) (string, error) {
	return "ShortURL", nil
}

func (svc *Service) GetLongURL(shortURL string) (string, error) {
	return "LongURL", nil
}
