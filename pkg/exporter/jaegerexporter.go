package exporter

import (
	"fmt"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewJaegerExporter() (sdktrace.SpanExporter, error) {
	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		fmt.Println(err)
		return (sdktrace.SpanExporter)(nil), err
	}
	return jaegerExporter, nil
}
