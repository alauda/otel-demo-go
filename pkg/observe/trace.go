package observe

import (
	"fmt"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"otel-demo-go/pkg/exporter"
)

func InitTracerProvider(serviceName string) *sdktrace.TracerProvider {
	jaegerExporter, err := exporter.NewJaegerExporter()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	tp := NewTraceProvider(serviceName, jaegerExporter)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
	return tp
}

func NewTraceProvider(serviceName string, exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
