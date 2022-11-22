package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func TestAddRouter(t *testing.T) {
	// 构建测试用例
	testcase := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/user"},
		{http.MethodGet, "/user/home"},
		{http.MethodGet, "/order/detail"},
		{http.MethodPost, "/order/create"},
		{http.MethodPost, "/login"},

		{http.MethodGet, "/order/*"},
		{http.MethodGet, "/*"},
		{http.MethodGet, "/*/*"},
		{http.MethodGet, "/*/abc"},
		{http.MethodGet, "/*/abc/*"},

		{http.MethodGet, "/order/detail/:id"},
		{http.MethodGet, "/order/detail/:id/qwe"},

		{http.MethodGet, "/user/:id(^[0-9]+$)"},
	}

	// 调用api
	r := newRouter()
	var mockHandler HandleFunc = func(c *Context) {}
	for _, c := range testcase {
		r.addRoute(c.method, c.path, mockHandler)
	}

	// 创建真实结果
	reg, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		t.Error("regex error")
	}
	want := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:       "/",
				handleFunc: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:       "user",
						handleFunc: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:       "home",
								handleFunc: mockHandler,
								children:   map[string]*node{},
							},
						},
						pathParam: &node{
							path:       ":id",
							handleFunc: mockHandler,
							children:   map[string]*node{},
							regExpr:    reg,
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:       "detail",
								handleFunc: mockHandler,
								children:   map[string]*node{},
								pathParam: &node{
									path:       ":id",
									handleFunc: mockHandler,
									children: map[string]*node{
										"qwe": {
											path:       "qwe",
											handleFunc: mockHandler,
											children:   map[string]*node{},
										},
									},
								},
							},
						},
						wildcard: &node{
							path:       "*",
							handleFunc: mockHandler,
							children:   map[string]*node{},
						},
					},
				},
				wildcard: &node{
					path:       "*",
					handleFunc: mockHandler,
					children: map[string]*node{
						"abc": {
							path:       "abc",
							handleFunc: mockHandler,
							children:   map[string]*node{},
							wildcard: &node{
								path:       "*",
								handleFunc: mockHandler,
								children:   map[string]*node{},
							},
						},
					},
					wildcard: &node{
						path:       "*",
						handleFunc: mockHandler,
						children:   map[string]*node{},
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:       "create",
								handleFunc: mockHandler,
								children:   map[string]*node{},
							},
						},
					},
					"login": {
						path:       "login",
						handleFunc: mockHandler,
						children:   map[string]*node{},
					},
				},
			},
		},
	}

	// 比较相等
	msg, equal := r.equals(want)
	assert.True(t, equal, msg)

	// 测试无效path
	r = newRouter()
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "user", mockHandler)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/user/", mockHandler)
	})
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a//b", mockHandler)
	})

	// 测试重复注册
	r = newRouter()

	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})

	r.addRoute(http.MethodGet, "/a/b", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/a/b", mockHandler)
	})

	// 测试conflict
	// 测试重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/user/*", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/user/:id", mockHandler)
	})

	r = newRouter()
	r.addRoute(http.MethodGet, "/user/:id", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/user/*", mockHandler)
	})

	r = newRouter()
	r.addRoute(http.MethodGet, "/user/:id", mockHandler)
	assert.Panics(t, func() {
		r.addRoute(http.MethodGet, "/user/:detail", mockHandler)
	})
}

func (r *router) equals(y *router) (string, bool) {
	for method, root1 := range r.trees {
		root2, ok := y.trees[method]
		if !ok {
			return fmt.Sprintf("%s error", method), false
		}
		if msg, equal := root1.equals(root2); !equal {
			return msg, false
		}

	}
	return "", true
}

func (n *node) equals(y *node) (string, bool) {
	if n.path != y.path {
		return "path error", false
	}
	if len(n.children) != len(y.children) {
		return "children len not equal", false
	}
	nHandler := reflect.ValueOf(n.handleFunc)
	yHandler := reflect.ValueOf(y.handleFunc)
	if nHandler != yHandler {
		return "handler error", false
	}

	if n.wildcard != nil && y.wildcard != nil {
		msg, ok := n.wildcard.equals(y.wildcard)
		if !ok {
			return msg, false
		}
	} else if n.wildcard != nil || y.wildcard != nil {
		return "wildcard error", false
	}

	if n.pathParam != nil && y.pathParam != nil {
		msg, ok := n.pathParam.equals(y.pathParam)
		if !ok {
			return msg, false
		}
		if n.pathParam.regExpr != nil && y.pathParam.regExpr != nil {
			if n.pathParam.regExpr.String() != y.pathParam.regExpr.String() {
				return "regex not equal", false
			}
		} else if n.pathParam.regExpr != nil || y.pathParam.regExpr != nil {
			return "regex number not equal", false
		}

	} else if n.pathParam != nil || y.pathParam != nil {
		return "pathParam error", false
	}

	for k, root1 := range n.children {
		root2, ok := y.children[k]
		if !ok {
			return "prefix error", false
		}
		if msg, equal := root1.equals(root2); !equal {
			return msg, false
		}
	}

	return "", true
}

func TestFindRouter(t *testing.T) {
	routers := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/user"},
		{http.MethodGet, "/user/home"},
		{http.MethodGet, "/order/detail"},
		{http.MethodPost, "/order/create"},
		{http.MethodPost, "/login"},

		{http.MethodGet, "/order/*"},
		{http.MethodGet, "/*"},
		{http.MethodGet, "/*/*"},
		{http.MethodGet, "/*/abc"},
		{http.MethodGet, "/*/abc/*"},

		{http.MethodGet, "/order/detail/:id"},
		{http.MethodGet, "/match/:pid/:mid"},

		{http.MethodGet, "/a/b/*"},

		{http.MethodGet, "/user/:id(^[0-9]+$)"},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(c *Context) {}
	for _, c := range routers {
		r.addRoute(c.method, c.path, mockHandler)
	}

	//构建testcase
	testcase := []struct {
		name           string
		method         string
		path           string
		wantNode       *node
		wantExist      bool
		wantPathParams map[string]string
	}{
		{
			name:      "get:/",
			method:    http.MethodGet,
			path:      "/",
			wantNode:  r.trees[http.MethodGet],
			wantExist: true,
		},
		{
			name:      "get:/user",
			method:    http.MethodGet,
			path:      "/user",
			wantNode:  r.trees[http.MethodGet].children["user"],
			wantExist: true,
		},
		{
			name:      "get:/user/home",
			method:    http.MethodGet,
			path:      "/user/home",
			wantNode:  r.trees[http.MethodGet].children["user"].children["home"],
			wantExist: true,
		},
		{
			name:      "get:/order/detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantNode:  r.trees[http.MethodGet].children["order"].children["detail"],
			wantExist: true,
		},
		{
			name:      "post:/order/create",
			method:    http.MethodPost,
			path:      "/order/create",
			wantNode:  r.trees[http.MethodPost].children["order"].children["create"],
			wantExist: true,
		},
		{
			name:      "post:/login",
			method:    http.MethodPost,
			path:      "/login",
			wantNode:  r.trees[http.MethodPost].children["login"],
			wantExist: true,
		},
		{
			name:      "method not found",
			method:    http.MethodOptions,
			wantExist: false,
		},
		//{
		//	name:      "path not found",
		//	method:    http.MethodGet,
		//	path:      "/q/w/e",
		//	wantExist: false,
		//},
		{
			name:      "don't have handlefunc",
			method:    http.MethodGet,
			path:      "/order",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].children["order"],
		},
		{
			name:      "/order/*",
			method:    http.MethodGet,
			path:      "/order/aaaaaa",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].children["order"].wildcard,
		},
		{
			name:      "/*",
			method:    http.MethodGet,
			path:      "/ppp",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].wildcard,
		},
		{
			name:      "/*/*",
			method:    http.MethodGet,
			path:      "/ppp/nnn",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].wildcard.wildcard,
		},
		{
			name:      "/*/abc",
			method:    http.MethodGet,
			path:      "/b/abc",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].wildcard.children["abc"],
		},
		{
			name:      "/*/abc/*",
			method:    http.MethodGet,
			path:      "/q/abc/m",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].wildcard.children["abc"].wildcard,
		},
		//{
		//	name:      "/user/abc/home-/user/*/*",
		//	method:    http.MethodGet,
		//	path:      "/user/abc/detail",
		//	wantExist: true,
		//	wantNode:  r.trees[http.MethodGet].children["user"].wildcard.wildcard,
		//},
		{
			name:           "/order/detail/:id",
			method:         http.MethodGet,
			path:           "/order/detail/123",
			wantExist:      true,
			wantNode:       r.trees[http.MethodGet].children["order"].children["detail"].pathParam,
			wantPathParams: map[string]string{"id": "123"},
		},
		{
			name:           "/match/:pid/:mid",
			method:         http.MethodGet,
			path:           "/match/123/456",
			wantExist:      true,
			wantNode:       r.trees[http.MethodGet].children["match"].pathParam.pathParam,
			wantPathParams: map[string]string{"pid": "123", "mid": "456"},
		},
		{
			name:      "/a/b/* = /a/b/c/d/f/e",
			method:    http.MethodGet,
			path:      "/a/b/c/d/e/f",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].children["a"].children["b"].wildcard,
		},
		{
			name:      "/user/:id(^[0-9]+$)",
			method:    http.MethodGet,
			path:      "/user/123",
			wantExist: true,
			wantNode:  r.trees[http.MethodGet].children["user"].pathParam,
		},
		{
			name:      "/user/:id(^[0-9]+$) - no match",
			method:    http.MethodGet,
			path:      "/user/qwe",
			wantExist: false,
		},
	}

	// 调用api 拿到结果
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			match, exist := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantExist, exist)
			if !exist || (tc.wantExist == false && exist == true) {
				return
			}
			msg, ok := tc.wantNode.equals(match.node)
			assert.True(t, ok, msg)
			if tc.wantPathParams != nil {
				equal := reflect.DeepEqual(match.params, tc.wantPathParams)
				assert.True(t, equal, "path param error")
			}
		})
	}
}

func TestFetchRegex(t *testing.T) {
	seg, regex := fetchRegexp(":id(^[0-9]+$)")
	fmt.Println(seg)
	fmt.Println(regex.String())
}

// 测试路由中间件
func TestRouteMiddleware(t *testing.T) {

	actual := []string{}

	var md1 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md1")
			next(c)
		}
	}
	var md2 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md2")
			next(c)
		}
	}
	var md3 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md3")
			next(c)
		}
	}
	var md4 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md4")
			next(c)
		}
	}
	var md5 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md5")
			next(c)
		}
	}
	var md6 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md6")
			next(c)
		}
	}
	var md7 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md7")
			next(c)
		}
	}
	var md8 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md8")
			next(c)
		}
	}
	var md9 Middleware = func(next HandleFunc) HandleFunc {
		return func(c *Context) {
			actual = append(actual, "md9")
			next(c)
		}
	}

	r := newRouter()
	var mockHandleFunc HandleFunc = func(c *Context) {}
	r.addRoute(http.MethodGet, "/*/b", mockHandleFunc)
	r.addRoute(http.MethodGet, "/a/d", mockHandleFunc)
	r.addRoute(http.MethodGet, "/a/b/c", mockHandleFunc)
	r.addRoute(http.MethodGet, "/a/:id/c", mockHandleFunc)

	_ = r.addMiddlewares(http.MethodGet, "/", md1)
	_ = r.addMiddlewares(http.MethodGet, "/a", md2)
	_ = r.addMiddlewares(http.MethodGet, "/*", md3)
	_ = r.addMiddlewares(http.MethodGet, "/a/:id", md4)
	_ = r.addMiddlewares(http.MethodGet, "/a/b", md5)
	_ = r.addMiddlewares(http.MethodGet, "/a/d", md6)
	_ = r.addMiddlewares(http.MethodGet, "/*/b", md7)
	_ = r.addMiddlewares(http.MethodGet, "/a/:id/c", md8)
	_ = r.addMiddlewares(http.MethodGet, "/a/b/c", md9)

	testcase := []struct {
		name   string
		method string
		path   string
		want   []string
	}{
		{
			name:   "1",
			method: http.MethodGet,
			path:   "/a/b/c",
			want:   []string{"md1", "md3", "md2", "md7", "md4", "md5", "md8", "md9"},
		},
	}

	for _, tt := range testcase {
		t.Run(tt.name, func(t *testing.T) {
			matched, ok := r.findRoute(tt.method, tt.path)
			assert.True(t, ok)

			cur := matched.handleFunc
			for i := len(matched.matchedMiddlewares) - 1; i >= 0; i-- {
				cur = matched.matchedMiddlewares[i](cur)
			}
			cur(&Context{})
			equal := reflect.DeepEqual(tt.want, actual)
			assert.True(t, equal, "mdls not equal")
		})
	}
}
