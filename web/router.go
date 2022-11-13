package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type node struct {
	path       string
	handleFunc HandleFunc
	children   map[string]*node
	wildcard   *node
	pathParam  *node
	regExpr    *regexp.Regexp
}

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
			panic(fmt.Sprintf("path %s already registed", path))
		}
		root.handleFunc = handleFunc
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
}

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

// getChildOrCreate: get child node if existed, otherwise create and return
func (n *node) getChildOrCreate(seg string) *node {
	if seg[0] == ':' {
		seg, regex := fetchRegexp(seg)
		if n.pathParam != nil {
			if n.pathParam.path != seg || n.pathParam.regExpr.String() != regex.String() {
				panic(fmt.Sprintf("'%s' is conflict with existed path param '%s'", seg, n.pathParam.path))
			}
			return n.pathParam
		}
		// 下面是不存在pathparam 想新注册的情况
		if n.wildcard != nil {
			panic(fmt.Sprintf("'%s' is conflict with existed wildcard '%s'", seg, n.wildcard.path))
		}
		n.pathParam = &node{
			path:     seg,
			children: map[string]*node{},
			regExpr:  regex,
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
			}
		}
		return n.wildcard
	}

	if _, ok := n.children[seg]; !ok {
		n.children[seg] = &node{
			path:     seg,
			children: map[string]*node{},
		}
	}

	return n.children[seg]
}

// isValidPath: check path is valid
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

// getRootOrCreate: get root node for one http method if existed, otherwise create and return
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

type matchInfo struct {
	*node
	params map[string]string
}

// findRoute: get node according to http method and url path
func (r *router) findRoute(httpMethod, path string) (*matchInfo, bool) {
	root, ok := r.trees[httpMethod]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{node: root}, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	return root.getMatchInfo(segs)
}

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

func (n *node) childOf(seg string) (*node, bool) {
	if child, ok := n.children[seg]; ok {
		return child, true
	}
	return n.getNotNilWildcardOrPathparam()
}

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

// ============================================================
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
