package monitoring

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var Tracer = otel.Tracer("rest.http")

// InitTracer sets up the OpenTelemetry tracing pipeline.
// It returns a function that can be used to shut down the tracer.
// The function should be called when the application is shutting down.
func InitTracer() (func(context.Context) error, error) {
	ctx := context.Background()

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithEndpoint(os.Getenv("OTEL_COLLECTOR")), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	hostname, _ := os.Hostname()
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.HostName(hostname),
		semconv.ServiceName("rest"),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// OtelMiddleware is a middleware that creates a new span for each request
// and passes it along to the next handler.
// The span is added to the context of the request.
// The span is closed when the handler returns.
func OtelMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := Tracer.Start(c.Request().Context(), c.Path(), trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		span.SetAttributes(
			attribute.KeyValue{Key: "env", Value: attribute.StringValue(os.Getenv("ENVIRONMENT"))},
			attribute.KeyValue{Key: "db", Value: attribute.StringValue(os.Getenv("POSTGRES_HOST"))},
		)

		// Pass the context with the span along to the next handler
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}

func NewSpan(ctx context.Context, name string, statusCode codes.Code, statusDescription string, kv ...attribute.KeyValue) trace.Span {
	_, span := Tracer.Start(ctx, name)

	span.SetStatus(statusCode, statusDescription)
	span.SetAttributes(kv...)

	return span
}

func NewSpanWithParent(ctx context.Context, name string) (trace.Span, trace.Span) {
	parent := trace.SpanFromContext(ctx)
	_, span := Tracer.Start(ctx, name)

	return parent, span
}
