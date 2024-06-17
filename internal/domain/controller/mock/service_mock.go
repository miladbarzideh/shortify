package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
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
