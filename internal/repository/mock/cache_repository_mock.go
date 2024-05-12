package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/miladbarzideh/shortify/internal/model"
)

type CacheRepository struct {
	mock.Mock
}

func (m *CacheRepository) Set(ctx context.Context, url model.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *CacheRepository) Get(ctx context.Context, shortCode string) (model.URL, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(model.URL), args.Error(1)
}

func (m *CacheRepository) BuildKeyWithPrefix(url string) string {
	args := m.Called(url)
	return args.String(0)
}
