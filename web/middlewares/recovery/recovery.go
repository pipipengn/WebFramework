package recovery

import "WebFramework/web"

type MiddlewareBuilder struct {
	statusCode int
	data       []byte
	log        func(c *web.Context)
}

type Options struct {
	StatusCode int
	Data       []byte
	Log        func(c *web.Context)
}

func NewBuilder(opt Options) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		statusCode: opt.StatusCode,
		data:       opt.Data,
		log:        opt.Log,
	}
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(c *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					c.RespData = m.data
					c.RespStatusCode = m.statusCode
					m.log(c)
				}
			}()
			next(c)
		}
	}
}
