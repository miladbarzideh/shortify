package mock

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type Service struct {
	mock.Mock
}

func (m *Service) CreateShortURL(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *Service) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func (m *Service) BuildShortURL(shortCode string) string {
	return fmt.Sprintf("localhost:8513/api/v1/urls/%s", shortCode)
}

func (m *Service) CreateShortURLWithRetries(longURL string, shortCode string) error {
	args := m.Called(longURL, shortCode)
	return args.Error(0)
}
