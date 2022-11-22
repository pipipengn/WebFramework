package opentelemetry

import (
	"WebFramework/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const defaultInstrumentationName = "WebFramework/web/widdlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(c *web.Context) {
			// 和客户端的tracer连在一起,上游调用
			ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
			ctx, span := m.Tracer.Start(ctx, "unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", c.Request.Method))
			span.SetAttributes(attribute.String("http.url", c.Request.URL.String()))
			span.SetAttributes(attribute.String("http.schema", c.Request.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", c.Request.Host))
			c.Request = c.Request.WithContext(ctx)
			next(c)
			span.SetName(c.MatchedRoute)
			span.SetAttributes(attribute.Int("http.status", c.RespStatusCode))
		}
	}
}
