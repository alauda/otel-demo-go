package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"otel-demo-go/pkg/observe"
	"otel-demo-go/pkg/provider"
	"time"
)

func main() {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(otelgin.Middleware(provider.ServiceName))

	tp := observe.InitTracerProvider(provider.ServiceName)
	observe.InitMetricsExporter(provider.ServiceName, ":9464")

	if err := provider.InitDBTracing(); err != nil {
		fmt.Println(err)
		return
	}

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

	app.GET("/hello", provider.Hello)
	app.GET("/stu/:id", provider.GetStudentByID)
	app.Run(":8080")
}
