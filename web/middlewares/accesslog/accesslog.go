package accesslog

import (
	"WebFramework/web"
	"encoding/json"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) LogFunc(logfunc func(log string)) *MiddlewareBuilder {
	m.logFunc = logfunc
	return m
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(c *web.Context) {
			defer func() {
				l := acesslog{
					Host:       c.Request.Host,
					Router:     c.MatchedRoute,
					HttpMethod: c.Request.Method,
					Path:       c.Request.URL.Path,
				}
				bytes, _ := json.Marshal(l)
				m.logFunc(string(bytes))
			}()
			next(c)
		}
	}
}

type acesslog struct {
	Host       string `json:"host,omitempty"`
	Router     string `json:"router,omitempty"`
	HttpMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
