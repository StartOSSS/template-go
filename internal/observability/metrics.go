package observability

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Metrics exposes a Prometheus handler and shutdown hook for the meter provider.
type Metrics struct {
	Handler  http.Handler
	Shutdown func(context.Context) error
}

func SetupMetrics() (*Metrics, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(provider)
	return &Metrics{
		Handler:  promhttp.Handler(),
		Shutdown: provider.Shutdown,
	}, nil
}
