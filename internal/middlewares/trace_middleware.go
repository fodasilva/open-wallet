package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		spanName := c.FullPath()
		if spanName == "" {
			spanName = c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(c.Request.Method),
				semconv.HTTPURL(c.Request.URL.String()),
				semconv.HTTPRoute(c.FullPath()),
			),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(semconv.HTTPStatusCode(c.Writer.Status()))
	}
}
