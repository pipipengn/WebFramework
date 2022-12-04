package orm

import (
	"WebFramework/orm/internal/errs"
	"strings"
)

type Deleter[T any] struct {
	table   string
	where   []Predicate
	args    []any
	builder strings.Builder
	model   *model
}

func NewDeleter[T any]() *Deleter[T] {
	return &Deleter[T]{}
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.builder.WriteString("DELETE FROM ")
	m, err := parseModel(new(T))
	if err != nil {
		return nil, err
	}
	d.model = m

	tableName := m.tableName
	if d.table != "" {
		tableName = d.table
	}
	d.builder.WriteString("`")
	d.builder.WriteString(tableName)
	d.builder.WriteString("`")

	if len(d.where) > 0 {
		d.builder.WriteString(" WHERE ")
		p := d.where[0]
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		if err := d.buildExpression(p); err != nil {
			return nil, err
		}
	}

	return &Query{
		SQL:  d.builder.String() + ";",
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case value:
		d.builder.WriteString("?")
		d.addArg(exp.val)
	case Column:
		d.builder.WriteString("`")
		field, ok := d.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknowField(exp.name)
		}
		d.builder.WriteString(field.colName)
		d.builder.WriteString("`")
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			d.builder.WriteString("(")
		}
		if err := d.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			d.builder.WriteString(")")
		}

		d.builder.WriteString(exp.op.String())

		_, ok = exp.right.(Predicate)
		if ok {
			d.builder.WriteString("(")
		}
		if err := d.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			d.builder.WriteString(")")
		}
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}

	return nil
}

func (d *Deleter[T]) addArg(arg any) {
	if d.args == nil {
		d.args = make([]any, 0, 8)
	}
	d.args = append(d.args, arg)
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.where = append(d.where, predicates...)
	return d
}
