package web

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// 路由树
type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type node struct {
	path        string
	fullPath    string
	handleFunc  HandleFunc
	children    map[string]*node
	wildcard    *node
	pathParam   *node
	regExpr     *regexp.Regexp
	middlewares []Middleware
}

// =========================================================================================================

// addRoute: register url to router
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
// - 可以注册 /user/:a/:b   /user/:a
// - 不能注册 /user/:a/:b   /user/:c
func (r *router) addRoute(httpMethod, path string, handleFunc HandleFunc) {
	isValidPath(path)
	root := r.getRootOrCreate(httpMethod)

	if path == "/" {
		if root.handleFunc != nil {
			panic(fmt.Sprintf("'%s' conflict with existed path", path))
		}
		root.handleFunc = handleFunc
		root.fullPath = "/"
		fmt.Println("/")
		return
	}

	segs := strings.Split(path, "/")[1:]
	cur := root
	for _, seg := range segs {
		if seg == "" {
			panic("invalid path")
		}
		cur = cur.getChildOrCreate(seg)
	}

	if cur.handleFunc != nil {
		panic(fmt.Sprintf("'%s' conflict with existed path", path))
	}
	cur.handleFunc = handleFunc
	fmt.Println(cur.fullPath)
}

// addRoute： get child node if existed, otherwise create and return
func (n *node) getChildOrCreate(seg string) *node {
	if seg[0] == ':' {
		// 提取正则
		seg, regex := fetchRegexp(seg)
		if n.pathParam != nil {
			if n.pathParam.path != seg ||
				(regex != nil && n.pathParam.regExpr != nil && n.pathParam.regExpr.String() != regex.String()) {
				panic(fmt.Sprintf("'%s' is conflict with existed path param '%s'", seg, n.pathParam.path))
			}
			return n.pathParam
		}
		// 下面是不存在pathparam 想新注册的情况
		if n.wildcard != nil {
			panic(fmt.Sprintf("'%s' is conflict with existed wildcard '%s'", seg, n.wildcard.path))
		}
		newseg := seg
		if regex != nil {
			newseg += fmt.Sprintf("(%s)", regex.String())
		}
		n.pathParam = &node{
			path:     seg,
			children: map[string]*node{},
			regExpr:  regex,
			fullPath: createFullPath(n, newseg),
		}
		return n.pathParam
	}

	if seg == "*" {
		if n.pathParam != nil {
			panic(fmt.Sprintf("'%s' is conflict with existed path param", seg))
		}
		if n.wildcard == nil {
			n.wildcard = &node{
				path:     "*",
				children: map[string]*node{},
				fullPath: createFullPath(n, seg),
			}
		}
		return n.wildcard
	}

	if _, ok := n.children[seg]; !ok {
		n.children[seg] = &node{
			path:     seg,
			children: map[string]*node{},
			fullPath: createFullPath(n, seg),
		}
	}

	return n.children[seg]
}

// addRoute：提取用户注册的路由中的正则
func fetchRegexp(seg string) (string, *regexp.Regexp) {
	for i, r := range seg {
		if r != '(' {
			continue
		}
		if !strings.HasSuffix(seg, ")") {
			panic(fmt.Sprintf("regex format error for %s", seg))
		}
		regex, err := regexp.Compile(seg[i+1 : len(seg)-1])
		if err != nil {
			panic(fmt.Sprintf("regex format error for %s", seg))
		}
		return seg[:i], regex
	}
	return seg, nil
}

// addRoute：构建从root开始到当前节点到fullpath
func createFullPath(n *node, path string) string {
	if n.fullPath == "/" {
		return "/" + path
	}
	return n.fullPath + "/" + path
}

// addRoute： check path is valid
func isValidPath(path string) {
	if path == "" {
		panic("path cannot be empty")
	}
	if !strings.HasPrefix(path, "/") {
		panic("path must starts with /")
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		panic("path mustn't end with /")
	}
}

// addRoute： get root node for one http method if existed, otherwise create and return
func (r *router) getRootOrCreate(httpMethod string) *node {
	root, ok := r.trees[httpMethod]
	if !ok {
		root = &node{
			path:     "/",
			children: map[string]*node{},
		}
		r.trees[httpMethod] = root
	}
	return root
}

// ================================================================================================================

// 匹配到的结果
type matchInfo struct {
	*node
	params             map[string]string
	matchedMiddlewares []Middleware
}

// findRoute: get node according to http method and url path
func (r *router) findRoute(httpMethod, path string) (*matchInfo, bool) {
	root, ok := r.trees[httpMethod]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			node:               root,
			matchedMiddlewares: root.middlewares, // 给root添加路由中间件
		}, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	matched, ok := root.getMatchInfo(segs)
	if !ok {
		return nil, false
	}
	// 去匹配路由中间件
	matched.matchedMiddlewares = getMatchedMiddlewares(root, segs)

	return matched, true
}

// findRoute: 匹配路由中间件
func getMatchedMiddlewares(root *node, segs []string) []Middleware {
	res := []Middleware{}
	res = append(res, root.middlewares...) // root上的中间件单独加上
	queue := []*node{root}
	level := 0

	for len(queue) > 0 && level < len(segs) {
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]
			// 越具体越后调度，所以1.通配符 2.路径参数 3.精准路由
			// 通配符
			if cur.wildcard != nil {
				res = append(res, cur.wildcard.middlewares...)
				queue = append(queue, cur.wildcard)
			}
			// 路径参数
			if cur.pathParam != nil {
				res = append(res, cur.pathParam.middlewares...)
				queue = append(queue, cur.pathParam)
			}
			// 精准路由
			if v, ok := cur.children[segs[level]]; ok {
				res = append(res, v.middlewares...)
				queue = append(queue, v)
			}
		}
		level++
	}

	return res
}

// findRoute: getMatchInfo
func (n *node) getMatchInfo(segs []string) (*matchInfo, bool) {
	res := &matchInfo{params: map[string]string{}}
	cur := n
	for _, seg := range segs {
		child, ok := cur.childOf(seg)
		if !ok {
			return nil, false
		}
		// 路径参数
		if strings.HasPrefix(child.path, ":") {
			// 如果这个节点上有正则，就去验证一下是否匹配
			if child.regExpr != nil && !child.regExpr.MatchString(seg) {
				return nil, false
			}
			// 把url中的路径参数带出来
			res.params[child.path[1:]] = seg
		}
		// 支持末尾通配符匹配多段
		if child.path == "*" && isLeaf(child) && child.handleFunc != nil {
			res.node = child
			return res, true
		}
		cur = child
	}

	res.node = cur
	return res, true
}

// findRoute：获取当前节点的child
func (n *node) childOf(seg string) (*node, bool) {
	if child, ok := n.children[seg]; ok {
		return child, true
	}
	return n.getNotNilWildcardOrPathparam()
}

// findRoute：获取当前节点的pathparam或wildcard
func (n *node) getNotNilWildcardOrPathparam() (*node, bool) {
	if n.pathParam != nil && n.wildcard != nil {
		panic("conflict between path param and wildcard")
	}
	if n.pathParam != nil {
		return n.pathParam, true
	}
	if n.wildcard != nil {
		return n.wildcard, true
	}
	return nil, false
}

func isLeaf(child *node) bool {
	if len(child.children) == 0 && child.wildcard == nil && child.pathParam == nil {
		return true
	}
	return false
}

// ========================================
// getNodeV2: support backtracking matching
func (n *node) getNodeV2(segs []string) (*node, bool) {
	return n.getNodeV2DFS(append([]string{"/"}, segs...), 0)
}

func (n *node) getNodeV2DFS(segs []string, index int) (*node, bool) {
	if index == len(segs)-1 && (segs[index] == n.path || n.path == "*") {
		return n, true
	}
	if index == len(segs)-1 {
		return nil, false
	}

	children, ok := n.childOfV2(segs[index+1])
	if !ok {
		return nil, false
	}

	for _, child := range children {
		if node, ok := child.getNodeV2DFS(segs, index+1); ok {
			return node, true
		}
	}
	return nil, false
}

func (n *node) childOfV2(seg string) ([]*node, bool) {
	res := []*node{}

	if child, ok := n.children[seg]; ok {
		res = append(res, child)
	}
	if n.wildcard != nil {
		res = append(res, n.wildcard)
	}
	return res, len(res) > 0
}

// =========================================================================================================

// 往对应的节点添加middleware
func (r *router) addMiddlewares(httpMethod, path string, middlewares ...Middleware) error {
	matched, ok := r.findRoute(httpMethod, path)
	if !ok {
		return errors.New(fmt.Sprintf("[method:%s] [path:%s] not exist", httpMethod, path))
	}
	matched.middlewares = append(matched.middlewares, middlewares...)
	return nil
}
