package web

import (
	"fmt"
	"net"
	"net/http"
)

// 确保一定实现接口
var _ Server = &httpServer{}

// HandleFunc 定义业务处理函数
type HandleFunc func(c *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(httpMethod, path string, handleFunc HandleFunc)
	addMiddlewares(httpMethod, path string, middlewares ...Middleware) error
}

// 默认实现类
type httpServer struct {
	*router
	middlewares []Middleware
	log         func(msg string, args ...any)
}

func NewHttpServer(opts ...HttpServerOption) *httpServer {
	res := &httpServer{
		router: newRouter(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

type HttpServerOption func(server *httpServer)

// WithMiddleware 初始化server的时候可以添加中间件
func WithMiddleware(middlewares ...Middleware) HttpServerOption {
	return func(server *httpServer) {
		server.middlewares = middlewares
	}
}

// ServeHTTP 处理请求的入口
func (h *httpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := &Context{
		Request: request,
		Writer:  writer,
	}

	// 把中间件串起来
	cur := h.serve
	for i := len(h.middlewares) - 1; i >= 0; i-- {
		cur = h.middlewares[i](cur)
	}

	// 添加最前面的flush中间件
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			next(c)
			h.flushResp(c)
		}
	}
	start := m(cur)
	start(c)
}

// 路由匹配并开始执行业务逻辑
func (h *httpServer) serve(c *Context) {
	match, ok := h.findRoute(c.Request.Method, c.Request.URL.Path)
	if !ok || match.handleFunc == nil {
		c.RespStatusCode = 404
		c.RespData = []byte("Not Found")
		return
	}

	// 将匹配到到路由中间件串起来
	cur := match.handleFunc
	for i := len(match.matchedMiddlewares) - 1; i >= 0; i-- {
		cur = match.matchedMiddlewares[i](cur)
	}

	c.Params = match.params
	c.MatchedRoute = match.fullPath
	cur(c)
}

// flushResp 最后一次性往前端发数据
func (h *httpServer) flushResp(c *Context) {
	if c.RespStatusCode != 0 {
		c.Writer.WriteHeader(c.RespStatusCode)
	}
	if n, err := c.Writer.Write(c.RespData); err != nil || n != len(c.RespData) {
		h.log("response error: %v", err)
	}
}

// Start 启动server
func (h *httpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(l, h)
}

// MatchRoute 不启动server测试路由能否匹配上
func (h *httpServer) MatchRoute(route, path string) bool {
	r := newRouter()
	r.addRoute(http.MethodGet, route, func(c *Context) {})
	if _, ok := r.findRoute(http.MethodGet, path); !ok {
		return false
	}
	return true
}

// Use 添加中间件 - 在server上直接添加中间件
func (h *httpServer) Use(middlewares ...Middleware) {
	h.middlewares = append(h.middlewares, middlewares...)
}

// UseWithRoute 在路由树上添加中间件
func (h *httpServer) UseWithRoute(method, path string, middlewares ...Middleware) {
	if err := h.addMiddlewares(method, path, middlewares...); err != nil {
		panic(err.Error())
	}
}

// =======================================================================

func (h *httpServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *httpServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *httpServer) Put(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

func (h *httpServer) Patch(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPatch, path, handleFunc)
}

func (h *httpServer) Delete(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodDelete, path, handleFunc)
}

func (h *httpServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

func (h *httpServer) Head(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodHead, path, handleFunc)
}

func (h *httpServer) Trace(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodTrace, path, handleFunc)
}

func (h *httpServer) Connect(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodConnect, path, handleFunc)
}
