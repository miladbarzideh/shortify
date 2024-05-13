package infra

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
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

var NOOPTelemetry = Telemetry{
	TraceProvider: noop.NewTracerProvider(),
	MeterProvider: mnoop.NewMeterProvider(),
}

type shutdown interface {
	Shutdown(ctx context.Context) error
}

type Telemetry struct {
	TraceProvider trace.TracerProvider
	MeterProvider metric.MeterProvider
}

func NewTelemetry(cfg *Config) (*Telemetry, error) {
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	tracerProvider, err := newTraceProvider(cfg.Trace)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracerProvider)
	meterProvider, err := newMeterProvider()
	if err != nil {
		return nil, err
	}

	otel.SetMeterProvider(meterProvider)

	return &Telemetry{
		TraceProvider: tracerProvider,
		MeterProvider: meterProvider,
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

	return err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(cfg Trace) (traceProvider trace.TracerProvider, err error) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.JaegerHost), jaeger.WithAgentPort(cfg.JaegerPort)),
	)
	if err != nil {
		return
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNamespaceKey.String(cfg.ServiceNameSpaceKey),
			semconv.ServiceNameKey.String(cfg.ServiceNameKey),
		),
	)
	if err != nil {
		return
	}

	traceProvider = sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(sdkTrace.NewBatchSpanProcessor(exporter)),
		sdkTrace.WithResource(res),
	)

	return
}

func newMeterProvider() (metricProvider metric.MeterProvider, err error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return
	}

	metricProvider = sdkMetric.NewMeterProvider(
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			sdkMetric.WithInterval(1*time.Minute))),
	)

	return
}
