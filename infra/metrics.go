package infra

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric"
)

type Counter struct {
	counter metric.Int64Counter
}

func NewCounter(meter metric.Meter, name string) Counter {
	counter, err := meter.Int64Counter(name, metric.WithDescription(fmt.Sprintf("total number of %s", name)))
	if err != nil {
		panic(err)
	}

	return Counter{
		counter: counter,
	}
}

func (c Counter) Inc(ctx context.Context) {
	c.counter.Add(ctx, 1)
}

type Latency struct {
	histogram metric.Float64Histogram
}

func NewLatency(meter metric.Meter, name string) Latency {
	histogram, err := meter.Float64Histogram(name, metric.WithDescription(fmt.Sprintf("latency of %s", name)))
	if err != nil {
		panic(err)
	}

	return Latency{
		histogram: histogram,
	}
}

func (l Latency) Record(ctx context.Context, start time.Time) {
	l.histogram.Record(ctx, time.Since(start).Seconds())
}

type CacheStats struct {
	Hits   Counter
	Misses Counter
}

func NewCacheStats(meter metric.Meter) CacheStats {
	return CacheStats{
		Hits:   NewCounter(meter, "cache.hits"),
		Misses: NewCounter(meter, "cache.misses"),
	}
}
