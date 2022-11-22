package web

import (
	"fmt"
	"testing"
)

func TestServerE2E(t *testing.T) {
	s := NewHttpServer()
	s.middlewares = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(c *Context) {
				fmt.Println("第1个before")
				next(c)
				fmt.Println("第1个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(c *Context) {
				fmt.Println("第2个before")
				next(c)
				fmt.Println("第2个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(c *Context) {
				fmt.Println("第3个before")
				fmt.Println("第3个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(c *Context) {
				fmt.Println("第4个before")
				fmt.Println("第4个after")
			}
		},
	}
	s.ServeHTTP(nil, nil)
}
