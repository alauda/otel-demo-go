package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"otel-demo-go/pkg/httpclient"
	"otel-demo-go/pkg/observe"
	"time"
)

const ServiceName = "otel-demo-consumer"

var (
	tracer     = otel.Tracer(ServiceName)
	counter, _ = observe.NewCounter("request.counter", "this is a test counter", unit.Dimensionless)
)

func Hello(ctx *gin.Context) {
	otelCtx, span := tracer.Start(ctx.Request.Context(), "consumer.Hello")
	ctx.Request = ctx.Request.WithContext(otelCtx)
	fmt.Println(span.IsRecording())
	if span.IsRecording() {
		span.SetAttributes(attribute.Int("hello", 1))
		span.SetAttributes(
			semconv.HTTPMethodKey.String("Get"),
			attribute.Int("hello", 1),
		)
	}

	defer span.End()

	span.AddEvent("httpclient start", trace.WithAttributes(attribute.Int64("timestamp", time.Now().UnixMilli())))
	resp, err := httpclient.DefaultClient().Get(otelCtx, "http://otel-demo-provider-go:8080/hello")
	if err != nil {
		span.SetStatus(500, err.Error())
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(500, err.Error())
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}
	fmt.Println(span.SpanContext().TraceID())

	span.AddEvent("httpclient end", trace.WithAttributes(attribute.Int64("timestamp", time.Now().UnixMilli())))
	counter.Add(otelCtx, 1, attribute.String("service.name", ServiceName))

	resultMap := make(map[string]interface{}, 0)
	_ = json.Unmarshal(body, &resultMap)
	ctx.JSON(200, resultMap)
}

func GetStudentByID(ctx *gin.Context) {
	otelCtx, span := tracer.Start(ctx.Request.Context(), "consumer.GetStudentByID")
	ctx.Request = ctx.Request.WithContext(otelCtx)
	fmt.Println(span.IsRecording())
	if span.IsRecording() {
		span.SetAttributes(attribute.Int("hello", 1))
		span.SetAttributes(
			semconv.HTTPMethodKey.String("Get"),
			attribute.Int("hello", 1),
		)
	}

	defer span.End()

	span.AddEvent("db query start", trace.WithAttributes(attribute.Int64("timestamp", time.Now().UnixMilli())))
	resp, err := httpclient.DefaultClient().Get(otelCtx, "http://otel-demo-provider-go:8080/stu/"+ctx.Param("id"))
	if err != nil {
		span.SetStatus(500, err.Error())
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.SetStatus(500, err.Error())
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}
	fmt.Println(span.SpanContext().TraceID())

	span.AddEvent("db query end", trace.WithAttributes(attribute.Int64("timestamp", time.Now().UnixMilli())))
	counter.Add(otelCtx, 1, attribute.String("service.name", ServiceName))

	resultMap := make(map[string]interface{}, 0)
	_ = json.Unmarshal(body, &resultMap)
	ctx.JSON(200, resultMap)
}
