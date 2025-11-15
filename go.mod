module github.com/example/template-go

go 1.21

require (
	github.com/bufbuild/connect-go v1.10.0
	github.com/jackc/pgx/v5 v5.4.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.26.0
	go.opentelemetry.io/otel/exporters/prometheus v0.49.0
	go.opentelemetry.io/otel/sdk/metric v1.26.0
	go.opentelemetry.io/otel/sdk/resource v1.26.0
	go.opentelemetry.io/otel/sdk/trace v1.26.0
)

require (
	go.opentelemetry.io/otel/attribute v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
)
