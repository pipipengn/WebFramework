package web

import (
	"fmt"
	"testing"
)

func TestServerE2E(t *testing.T) {
	s := NewHttpServer()
	s.Get("/user", func(c *Context) {
		if _, err := c.Writer.Write([]byte("aaa")); err != nil {
			fmt.Println("error", err.Error())
		}
	})
	_ = s.Start(":8080")
}
