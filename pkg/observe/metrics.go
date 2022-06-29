package observe

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"net/http"
)

var meter = global.MeterProvider().Meter("")

func InitMetricsExporter(serviceName, addr string) {
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		fmt.Printf("failed to initialize prometheus exporter: %w\n", err)
		return
	}

	global.SetMeterProvider(exporter.MeterProvider())

	http.HandleFunc("/", exporter.ServeHTTP)
	go func() {
		_ = http.ListenAndServe(addr, nil)
	}()

	fmt.Println("Prometheus server running on " + addr)
}

func NewCounter(name, description string, unit unit.Unit) (syncint64.Counter, error) {
	return meter.SyncInt64().Counter(
		name,
		instrument.WithUnit(unit),
		instrument.WithDescription(description),
	)
}

func CounterObserver(serviceName string) {
	counter, _ := meter.AsyncInt64().Counter(
		"some.prefix.counter_observer",
		instrument.WithUnit("1"),
		instrument.WithDescription("TODO"),
	)

	var number int64
	if err := meter.RegisterCallback(
		[]instrument.Asynchronous{
			counter,
		},
		// SDK periodically calls this function to collect data.
		func(ctx context.Context) {
			number++
			counter.Observe(ctx, number, attribute.String("service.name", serviceName))
		},
	); err != nil {
		panic(err)
	}
}
