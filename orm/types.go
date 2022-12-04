package orm

import (
	"context"
	"database/sql"
)

// Querier 用于select
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMany(ctx context.Context) ([]*T, error)
}

// Executer 用于增删改
type Executer interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// QueryBuilder 用于构建sql
type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
