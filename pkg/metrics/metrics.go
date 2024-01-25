package metrics

import (
	"coffee-shop/config"
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

type Metrics interface {
	IncHits(ctx context.Context, status int, method, path string)
	ObserveResponseTime(ctx context.Context, status int, method, path string, observeTime float64)
}

type OtelMetrics struct {
	HitsTotal metric.Int64Counter
	Hits      metric.Int64Counter
	Times     metric.Int64Histogram
}

// Create metrics with address and name
func CreateMetrics(name string) (Metrics, error) {
	meterProvider := otel.GetMeterProvider().Meter("")

	var metr OtelMetrics

	metr.HitsTotal, _ = meterProvider.Int64Counter(
		name + "_hits_total",
	)

	metr.Hits, _ = meterProvider.Int64Counter(
		name + "_hits",
	)

	metr.Times, _ = meterProvider.Int64Histogram(
		name + "_times")
	return &metr, nil
}

// IncHits
func (metr *OtelMetrics) IncHits(ctx context.Context, status int, method, path string) {
	metr.HitsTotal.Add(ctx, 1)
	metr.Hits.Add(ctx, 1, metric.WithAttributes(attribute.Key(method).String(path)))
}

// Observer response time
func (metr *OtelMetrics) ObserveResponseTime(ctx context.Context, status int, method, path string, observeTime float64) {
	metr.Times.Record(ctx, int64(observeTime), metric.WithAttributes(attribute.Key(method).String(path)))
}

func InitMetricsProvider(cfg *config.Config, serviceName string) (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)

	if err != nil {
		return nil, err
	}

	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(cfg.Otel.Host),
	)

	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)

	otel.SetMeterProvider(meterProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := metricExporter.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}, nil

}
