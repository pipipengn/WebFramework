package web

import (
	"net"
	"net/http"
)

var _ Server = &httpServer{}

type HandleFunc func(c *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(httpMethod, path string, handleFunc HandleFunc)
}

type httpServer struct {
	*router
}

func NewHttpServer() *httpServer {
	return &httpServer{newRouter()}
}

// ServeHTTP 处理请求的入口
func (h *httpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := &Context{
		Request: request,
		Writer:  writer,
	}
	h.serve(c)
}

func (h *httpServer) serve(c *Context) {
	match, ok := h.findRoute(c.Request.Method, c.Request.URL.Path)
	//!(match.path == "*" && match.isLastWildcard)
	if !ok || match.handleFunc == nil {
		c.Writer.WriteHeader(404)
		_, _ = c.Writer.Write([]byte("Not Found"))
		return
	}

	c.Params = match.params
	match.handleFunc(c)
}

func (h *httpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// TODO 做一些初始化

	return http.Serve(l, h)
}

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
