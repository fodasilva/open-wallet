package infra

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitTracer() (*trace.TracerProvider, error) {
	ctx := context.Background()
	var exporter sdktrace.SpanExporter
	var err error

	switch AppConfig.Enviroment {
	case "dev":

		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithInsecure())
	case "prod":
		exporter, err = texporter.New(texporter.WithProjectID(*AppConfig.GcpProjectID))
	}

	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("open-wallet"),
			attribute.String("environment", AppConfig.Enviroment),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
	return tp, nil
}
