package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/miladbarzideh/shortify/internal/domain/model"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) Create(ctx context.Context, longURL string, shortCode string) (model.URL, error) {
	args := m.Called(ctx, longURL, shortCode)
	return args.Get(0).(model.URL), args.Error(1)
}

func (m *Repository) FindByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
	args := m.Called(shortCode)
	return args.Get(0).(model.URL), args.Error(1)
}
