package trace

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/felipe1496/open-wallet/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() (*trace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error

	switch utils.AppConfig.Enviroment {
	case "dev":
		ctx := context.Background()
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint("localhost:4318"),
			otlptracehttp.WithInsecure())
	case "prod":
		exporter, err = texporter.New(texporter.WithProjectID(*utils.AppConfig.GcpProjectID))
	}

	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
