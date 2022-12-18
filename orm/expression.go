package orm

// Expression 标记接口 代表表达式
type Expression interface {
	expr()
}

// RawExpr 原生sql表达式
type RawExpr struct {
	raw  string
	args []any
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{raw: expr, args: args}
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{left: r}
}

func (r RawExpr) selectable() {}
func (r RawExpr) expr()       {}
