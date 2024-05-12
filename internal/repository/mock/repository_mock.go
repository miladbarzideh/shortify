package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/miladbarzideh/shortify/internal/model"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) Create(longURL string, shortCode string) (model.URL, error) {
	args := m.Called(longURL, shortCode)
	return args.Get(0).(model.URL), args.Error(1)
}

func (m *Repository) FindByShortCode(shortURL string) (model.URL, error) {
	args := m.Called(shortURL)
	return args.Get(0).(model.URL), args.Error(1)
}
