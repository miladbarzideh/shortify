package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/model"
)

type URLRepository interface {
	Create(longURL string, shortCode string) (model.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (model.URL, error)
}

type Repository struct {
	logger *logrus.Logger
	db     *gorm.DB
	tracer trace.Tracer
}

func NewRepository(logger *logrus.Logger, db *gorm.DB, tracer trace.Tracer) URLRepository {
	return &Repository{
		logger: logger,
		db:     db,
		tracer: tracer,
	}
}

func (r Repository) Create(longURL string, shortCode string) (model.URL, error) {
	url := model.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	result := r.db.Create(&url)

	return url, result.Error
}

func (r Repository) FindByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
	_, span := r.tracer.Start(ctx, "urlRepo.find")
	defer span.End()
	var url model.URL
	result := r.db.Where("short_code = ?", shortCode).First(&url)

	return url, result.Error
}
