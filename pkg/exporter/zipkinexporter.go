package exporter

import (
	"fmt"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewZipkinExporter() (sdktrace.SpanExporter, error) {
	zipkinExporter, err := zipkin.New("")
	if err != nil {
		fmt.Println(err)
		return (sdktrace.SpanExporter)(nil), err
	}
	return zipkinExporter, nil
}
