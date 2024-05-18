package infra

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	mnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var NOOPTelemetry = &Telemetry{
	TraceProvider: noop.NewTracerProvider(),
	MeterProvider: mnoop.NewMeterProvider(),
}

type shutdown interface {
	Shutdown(ctx context.Context) error
}

type Telemetry struct {
	TraceProvider trace.TracerProvider
	MeterProvider metric.MeterProvider
	logger        *logrus.Logger
	server        *http.Server
}

func NewTelemetry(logger *logrus.Logger, cfg *Config) (*Telemetry, error) {
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNamespaceKey.String(cfg.Telemetry.ServiceNamespaceKey),
			semconv.ServiceNameKey.String(cfg.Telemetry.ServiceNameKey),
		),
	)
	if err != nil {
		return nil, err
	}

	tracerProvider, err := newTraceProvider(cfg.Telemetry.Trace, res)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	meterProvider, err := newMeterProvider(cfg.Telemetry, res)
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(meterProvider)
	var server *http.Server
	if cfg.Telemetry.Metric.Enabled {
		server = &http.Server{Addr: cfg.Telemetry.Metric.Address}
		server.Handler = http.DefaultServeMux
		http.Handle("/metrics", promhttp.Handler())
		go func() {
			if err = server.ListenAndServe(); err != nil {
				logger.Fatalf("Failed to start metrics server: %v", err)
			}
		}()
	}

	return &Telemetry{
		TraceProvider: tracerProvider,
		MeterProvider: meterProvider,
		logger:        logger,
		server:        server,
	}, nil
}

func (t Telemetry) Shutdown(ctx context.Context) error {
	var err error
	if t.TraceProvider != nil {
		tp, ok := t.TraceProvider.(shutdown)
		if !ok {
			err = errors.Join(err, tp.Shutdown(ctx))
		}
	}

	if t.MeterProvider != nil {
		mp, ok := t.MeterProvider.(shutdown)
		if !ok {
			err = errors.Join(err, mp.Shutdown(ctx))
		}
	}

	if t.server != nil {
		errors.Join(err, t.server.Shutdown(ctx))
	}

	return err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(cfg Trace, resource *resource.Resource) (traceProvider trace.TracerProvider, err error) {
	if !cfg.Enabled {
		return NOOPTelemetry.TraceProvider, nil
	}

	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.JaegerHost), jaeger.WithAgentPort(cfg.JaegerPort)),
	)
	if err != nil {
		return
	}

	traceProvider = sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(sdkTrace.NewBatchSpanProcessor(exporter)),
		sdkTrace.WithResource(resource),
	)

	return
}

func newMeterProvider(cfg Tele, resource *resource.Resource) (metricProvider metric.MeterProvider, err error) {
	if !cfg.Metric.Enabled {
		return NOOPTelemetry.MeterProvider, nil
	}

	exporter, err := prometheus.New(prometheus.WithNamespace(cfg.ServiceNamespaceKey))
	if err != nil {
		return
	}

	metricProvider = sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(resource),
		sdkMetric.WithReader(exporter),
	)

	return
}
