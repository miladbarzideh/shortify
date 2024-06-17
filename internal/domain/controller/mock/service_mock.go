package mock

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/miladbarzideh/shortify/internal/domain/model"
)

type Service struct {
	mock.Mock
}

func (m *Service) CreateShortURL(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *Service) GetLongURL(ctx context.Context, shortCode string) (string, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Error(1)
}

func (m *Service) BuildShortURL(shortCode string) string {
	return fmt.Sprintf("localhost:8513/api/v1/urls/%s", shortCode)
}

func (m *Service) CreateShortURLWithRetries(ctx context.Context, longURL string, shortCode string) (*model.URL, error) {
	args := m.Called(ctx, longURL, shortCode)
	return args.Get(0).(*model.URL), args.Error(1)
}
