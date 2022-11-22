package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Request        *http.Request
	Writer         http.ResponseWriter
	Params         map[string]string
	queryValues    url.Values
	MatchedRoute   string
	RespData       []byte
	RespStatusCode int
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("nil input")
	}
	if c.Request.Body == nil {
		return errors.New("nil body")
	}
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) StringValue {
	if err := c.Request.ParseForm(); err != nil {
		return StringValue{Err: err}
	}
	return StringValue{Val: c.Request.FormValue(key)}
}

func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Request.URL.Query()
	}
	if _, ok := c.queryValues[key]; !ok {
		return StringValue{Err: errors.New("query key does not exist")}
	}
	return StringValue{Val: c.queryValues.Get(key)}
}

func (c *Context) PathValue(key string) StringValue {
	if val, ok := c.Params[key]; ok {
		return StringValue{Val: val}
	}
	return StringValue{Err: errors.New("path params key does not exist")}
}

type StringValue struct {
	Val string
	Err error
}

func (s *StringValue) AsInt64() (int64, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	return strconv.ParseInt(s.Val, 10, 64)
}

func (c *Context) JSON(status int, v any) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(status)
	bytes, _ := json.Marshal(v)
	c.RespData = bytes
	c.RespStatusCode = status
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}
