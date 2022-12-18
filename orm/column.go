package orm

type Column struct {
	name string
}

func Col(name string) Column {
	return Column{name: name}
}

func valueOf(val any) Expression {
	switch v := val.(type) {
	case Expression:
		return v
	default:
		return value{val: v}
	}
}

func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(arg),
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valueOf(arg),
	}
}

func (c Column) Lte(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLte,
		right: valueOf(arg),
	}
}

func (c Column) Gt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: valueOf(arg),
	}
}

func (c Column) Gte(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGte,
		right: valueOf(arg),
	}
}

func (Column) expr()         {}
func (c Column) selectable() {}
