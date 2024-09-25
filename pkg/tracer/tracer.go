package tracer

import (
	"context"
	"fmt"

	"garantexGRPC/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Trace struct {
	Provider *trace.TracerProvider
}

func NewProvider(cfg configs.TraceConfig) (*Trace, error) {

	provider, err := InitTracerProvider(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return &Trace{
		Provider: provider,
	}, nil
}

func InitTracerProvider(ctx context.Context, cfg configs.TraceConfig) (*trace.TracerProvider, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.Name),
		),
	)
	if err != nil {
		return nil, err
	}
	provider := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(res))
	otel.SetTracerProvider(provider)
	return provider, nil
}
