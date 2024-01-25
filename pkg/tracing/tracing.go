package tracing

import (
	"coffee-shop/config"
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"google.golang.org/grpc"
)

func InitTracerProvider(cfg *config.Config, serviceName string) (func(), error) {
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

	traceClient := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.Otel.Host),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	traceExporter, err := otlptrace.New(ctx, traceClient)

	if err != nil {
		return nil, err
	}

	bsp := sdk_trace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdk_trace.NewTracerProvider(
		sdk_trace.WithSampler(sdk_trace.AlwaysSample()),
		sdk_trace.WithResource(res),
		sdk_trace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExporter.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}, nil
}
