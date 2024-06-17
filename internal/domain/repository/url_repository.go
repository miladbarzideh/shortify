package repository

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/miladbarzideh/shortify/internal/domain/model"
	"github.com/miladbarzideh/shortify/internal/infra"
)

type Repository struct {
	logger        *logrus.Logger
	db            *gorm.DB
	tracer        trace.Tracer
	createLatency infra.Latency
	getLatency    infra.Latency
}

func NewRepository(logger *logrus.Logger, db *gorm.DB, telemetry *infra.TelemetryProvider) *Repository {
	tracer := telemetry.TraceProvider.Tracer("urlRepo")
	meter := telemetry.MeterProvider.Meter("urlRepo")
	createLatency := infra.NewLatency(meter, "db.create")
	getLatency := infra.NewLatency(meter, "db.get")

	return &Repository{
		logger:        logger,
		db:            db,
		tracer:        tracer,
		createLatency: createLatency,
		getLatency:    getLatency,
	}
}

func (r Repository) Create(ctx context.Context, url *model.URL) error {
	start := time.Now()
	result := r.db.Create(url)
	if result.Error != nil {
		return result.Error
	}

	r.createLatency.Record(ctx, start)

	return nil
}

func (r Repository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	start := time.Now()
	_, span := r.tracer.Start(ctx, "urlRepo.find")
	defer span.End()
	var url model.URL
	result := r.db.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return nil, result.Error
	}

	r.getLatency.Record(ctx, start)

	return &url, result.Error
}
