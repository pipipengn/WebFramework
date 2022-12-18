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

type value struct {
	val any
}

func (Predicate) expr() {}
func (value) expr()     {}
