package accesslog

import (
	"WebFramework/web"
	"fmt"
	"testing"
)

func TestBuilderE2E(t *testing.T) {
	middleware := NewBuilder().LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	s := web.NewHttpServer(web.WithMiddleware(middleware))
	s.Get("/a/b/*", func(c *web.Context) {
		fmt.Println("业务逻辑")
	})

	_ = s.Start(":8080")
}
