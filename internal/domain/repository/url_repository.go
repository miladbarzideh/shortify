package repository

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	infra2 "github.com/miladbarzideh/shortify/internal/infra"
)

type URLRepository interface {
	Create(ctx context.Context, longURL string, shortCode string) (model.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (model.URL, error)
}

type Repository struct {
	logger        *logrus.Logger
	db            *gorm.DB
	tracer        trace.Tracer
	createLatency infra2.Latency
	getLatency    infra2.Latency
}

func NewRepository(logger *logrus.Logger, db *gorm.DB, telemetry *infra2.TelemetryProvider) URLRepository {
	tracer := telemetry.TraceProvider.Tracer("urlRepo")
	meter := telemetry.MeterProvider.Meter("urlRepo")
	createLatency := infra2.NewLatency(meter, "db.create")
	getLatency := infra2.NewLatency(meter, "db.get")

	return &Repository{
		logger:        logger,
		db:            db,
		tracer:        tracer,
		createLatency: createLatency,
		getLatency:    getLatency,
	}
}

func (r Repository) Create(ctx context.Context, longURL string, shortCode string) (model.URL, error) {
	start := time.Now()
	url := model.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	result := r.db.Create(&url)
	if result.Error == nil {
		r.createLatency.Record(ctx, start)
	}

	return url, result.Error
}

func (r Repository) FindByShortCode(ctx context.Context, shortCode string) (model.URL, error) {
	start := time.Now()
	_, span := r.tracer.Start(ctx, "urlRepo.find")
	defer span.End()
	var url model.URL
	result := r.db.Where("short_code = ?", shortCode).First(&url)
	if result.Error == nil {
		r.getLatency.Record(ctx, start)
	}

	return url, result.Error
}
