package web

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Params  map[string]string
}

func (c Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) String(s string) {
	_, _ = c.Writer.Write([]byte(s))
}

func (c *Context) JSON(v any) {
	bytes, _ := json.Marshal(v)
	_, _ = c.Writer.Write(bytes)
}
