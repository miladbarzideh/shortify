package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/miladbarzideh/shortify/internal/domain/model"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) Create(ctx context.Context, url *model.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *Repository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(*model.URL), args.Error(1)
}
