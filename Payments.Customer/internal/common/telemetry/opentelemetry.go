package opentel

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type OpenTel struct {
	ServiceName      string
	ServiceVersion   string
	ExporterEndpoint string
	tracerProvider   *sdktrace.TracerProvider
	propagator       propagation.TextMapPropagator
}

func NewOpenTel() *OpenTel {
	return &OpenTel{
		propagator: propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	}
}

func (ot *OpenTel) GetTracerProvider() *sdktrace.TracerProvider {
	return ot.tracerProvider
}

func (ot *OpenTel) GetPropagators() propagation.TextMapPropagator {
	return ot.propagator
}

func (ot *OpenTel) GetTracer() trace.Tracer {
	ctx := context.Background()

	// Configure the OTLP exporter
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(ot.ExporterEndpoint),
		otlptracehttp.WithInsecure(), // This is important to use HTTP instead of HTTPS
	}

	client := otlptracehttp.NewClient(opts...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	// Configure the resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(ot.ServiceName),
		semconv.ServiceVersionKey.String(ot.ServiceVersion),
	)

	// Configure the tracer provider
	ot.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(ot.tracerProvider)
	otel.SetTextMapPropagator(ot.propagator)

	return ot.tracerProvider.Tracer(ot.ServiceName)
}
