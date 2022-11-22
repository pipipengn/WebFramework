package errhandle

import "WebFramework/web"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) AddCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(c *web.Context) {
			next(c)
			if data, ok := m.resp[c.RespStatusCode]; ok {
				c.RespData = data
			}
		}
	}
}
