package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"otel-demo-go/pkg/consumer"
	"otel-demo-go/pkg/observe"
	"time"
)

func main() {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(otelgin.Middleware(consumer.ServiceName))

	tp := observe.InitTracerProvider(consumer.ServiceName)
	observe.InitMetricsExporter("request.counter", ":9464")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	app.GET("/hello", consumer.Hello)
	app.GET("/stu/:id", consumer.GetStudentByID)
	app.Run(":8081")
}
