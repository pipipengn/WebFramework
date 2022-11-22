package accesslog

import (
	"WebFramework/web"
	"fmt"
	"net/http"
	"testing"
)

func TestBuilder(t *testing.T) {
	middleware := NewBuilder().LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	s := web.NewHttpServer(web.WithMiddleware(middleware))
	s.Post("/a/b/*", func(c *web.Context) {
		fmt.Println("业务逻辑")
	})

	request, _ := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	request.Host = "localhost"
	s.ServeHTTP(nil, request)
}
