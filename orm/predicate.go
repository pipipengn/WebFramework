package orm

type op string

func (o op) String() string {
	return string(o)
}

const (
	opEq  = "="
	opLt  = "<"
	opGt  = ">"
	opLte = "<="
	opGte = ">="
	opNot = "NOT"
	opAnd = "AND"
	opOr  = "OR"
)

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

type Column struct {
	name string
}

func Col(name string) Column {
	return Column{name: name}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: value{val: arg},
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: value{val: arg},
	}
}

func (c Column) Lte(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLte,
		right: value{val: arg},
	}
}

func (c Column) Gt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: value{val: arg},
	}
}

func (c Column) Gte(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGte,
		right: value{val: arg},
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

// Expression 标记接口 代表表达式
type Expression interface {
	expr()
}

type value struct {
	val any
}

func (Predicate) expr() {}
func (Column) expr()    {}
func (value) expr()     {}
