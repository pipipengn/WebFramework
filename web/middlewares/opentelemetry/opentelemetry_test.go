package opentelemetry

import (
	"WebFramework/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"os"
	"testing"
	"time"
)

func TestOpenTeleMetry(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	s := web.NewHttpServer(web.WithMiddleware(NewBuilder().Build()))
	s.Get("/user", func(c *web.Context) {
		ctx, span := tracer.Start(c.Request.Context(), "first_layer")
		defer span.End()

		secondC, second := tracer.Start(ctx, "second_layer")
		time.Sleep(time.Second)

		_, third1 := tracer.Start(secondC, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		_, third2 := tracer.Start(secondC, "third_layer_2")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()

		_, first := tracer.Start(c.Request.Context(), "first_layer_1")
		defer first.End()

		time.Sleep(100 * time.Millisecond)
		c.JSON(202, User{Name: "qwe"})
	})

	initZipkin(t)

	_ = s.Start(":8080")
}

type User struct {
	Name string
}

func initZipkin(t *testing.T) {
	exporter, err := zipkin.New(
		"http://localhost:19411/api/v2/spans",
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		t.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
}

func initJeager(t *testing.T) {
	url := "http://localhost:14268/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		t.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exp),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 1),
		)),
	)

	otel.SetTracerProvider(tp)
}
