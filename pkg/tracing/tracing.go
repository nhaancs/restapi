package tracing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tp *tracesdk.TracerProvider

// Config required for the tracing
type Config struct {
	Name    string
	Address string
}

func Init(c Config) error {
	parsedAddr := strings.Split(c.Address, ":")
	if len(parsedAddr) != 2 {
		return errors.New("invalid address")
	}
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(parsedAddr[0]), jaeger.WithAgentPort(parsedAddr[1])))
	if err != nil {
		return fmt.Errorf("create Jaeger exporter: %v", err)
	}
	initTracer(exp, c.Name)
	return nil
}

func InitForTest() {
	initTracer(tracetest.NewInMemoryExporter(), "")
}

func initTracer(exp tracesdk.SpanExporter, serviceName string) {
	tp = tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3SingleHeader), b3.WithInjectEncoding(b3.B3MultipleHeader)))
}

func Close(ctx context.Context) error {
	if tp != nil {
		return tp.Shutdown(ctx)
	}
	return errors.New("nil TraceProvider")
}

// AddEvent adds events and returns a context having span with that event
func AddEvent(ctx context.Context, s string) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(s)
}

// StartSpan creates a span and a context.Context containing the newly-created span.
func StartSpan(ctx context.Context, name string) context.Context {
	tr := otel.Tracer("")
	start, _ := tr.Start(ctx, name)
	return start
}

// SpanID returns span identity associated with the given context
// return "" if there is no trace associated
func SpanID(ctx context.Context) string {
	s := trace.SpanFromContext(ctx)
	return s.SpanContext().SpanID().String()
}

// EndSpan finishes the span that is associated with the given context
func EndSpan(ctx context.Context) {
	if sp := trace.SpanFromContext(ctx); sp != nil {
		sp.End()
	}
}

// TraceID return the traceID from the given context
func TraceID(ctx context.Context) string {
	if sp := trace.SpanFromContext(ctx); sp != nil {
		return sp.SpanContext().TraceID().String()
	}
	return ""
}
