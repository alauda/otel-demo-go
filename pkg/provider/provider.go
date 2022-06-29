package provider

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/extra/redisotel/v9"
	"github.com/go-redis/redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ServiceName = "otel-demo-provider"

var (
	tracer = otel.Tracer(ServiceName)
	db     *gorm.DB
	rdb    *redis.Client
)

type Student struct {
	ID   int
	Name string
	Age  int
}

func InitDBTracing() error {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	rdb.AddHook(redisotel.NewTracingHook(redisotel.WithAttributes(semconv.NetPeerNameKey.String("postgres"), semconv.NetPeerPortKey.String("6379"))))

	dsn := "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = db.Use(otelgorm.NewPlugin())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Hello(ctx *gin.Context) {
	fmt.Println(ctx.Request.Header.Get("x-b3-traceid"))
	_, span := tracer.Start(ctx.Request.Context(), "provider.Hello")
	defer span.End()

	ctx.JSON(200, gin.H{
		"foo": "bar",
	})
}

func GetStudentByID(ctx *gin.Context) {
	fmt.Println(ctx.Request.Header.Get("x-b3-traceid"))
	otelCtx, span := tracer.Start(ctx.Request.Context(), "provider.GetStudentByID")
	defer span.End()

	var stu Student
	err := db.WithContext(otelCtx).Where("id = ?", ctx.Param("id")).First(&stu).Error
	if err != nil {
		fmt.Println(err)
		ctx.Status(500)
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}

	if err := rdb.Set(otelCtx, "first_value", "value_1", 0).Err(); err != nil {
		fmt.Println(err)
		ctx.Status(500)
		span.RecordError(err, trace.WithAttributes(attribute.String("errmsg", err.Error())))
		return
	}

	ctx.JSON(200, stu)
}
