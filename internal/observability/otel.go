package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Setup configures OpenTelemetry tracing exporting via OTLP HTTP. It returns a
// shutdown function that must be invoked during graceful termination.
func Setup(ctx context.Context, serviceName, version string) (func(context.Context) error, error) {
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		"",
		attribute.String("service.name", serviceName),
		attribute.String("service.version", version),
	))
	if err != nil {
		return nil, fmt.Errorf("resource: %w", err)
	}

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("otlptrace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
