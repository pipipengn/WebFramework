package orm

import (
	"WebFramework/orm/internal/errs"
	"context"
	"strings"
)

type Selector[T any] struct {
	table   string
	builder strings.Builder
	where   []Predicate
	args    []any
	model   *model
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMany(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Build() (*Query, error) {
	s.builder.WriteString("SELECT * FROM ")
	m, err := parseModel(new(T))
	if err != nil {
		return nil, err
	}
	s.model = m

	if s.table == "" {
		s.builder.WriteString("`")
		s.builder.WriteString(m.tableName)
		s.builder.WriteString("`")
	} else {
		s.builder.WriteString(s.table)
	}

	if len(s.where) > 0 {
		s.builder.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	return &Query{
		SQL:  s.builder.String() + ";",
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
		return nil
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.builder.WriteString("(")
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.builder.WriteString(")")
		}

		s.builder.WriteString(exp.op.String())

		_, ok = exp.right.(Predicate)
		if ok {
			s.builder.WriteString("(")
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.builder.WriteString(")")
		}
	case Column:
		s.builder.WriteString("`")
		field, ok := s.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknowField(exp.name)
		}
		s.builder.WriteString(field.colName)
		s.builder.WriteString("`")
	case value:
		s.builder.WriteString("?")
		s.addArg(exp.val)
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (s *Selector[T]) addArg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 8)
	}
	s.args = append(s.args, val)
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(p ...Predicate) *Selector[T] {
	s.where = append(s.where, p...)
	return s
}
