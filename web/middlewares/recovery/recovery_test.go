package recovery

import (
	"WebFramework/web"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewBuilder(Options{
		StatusCode: 500,
		Data:       []byte("你panic了"),
		Log: func(c *web.Context) {

		},
	})
	s := web.NewHttpServer(web.WithMiddleware(builder.Build()))
	s.Get("/user", func(c *web.Context) {
		panic("panic")
	})
	_ = s.Start(":8080")
}
