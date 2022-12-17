package valuer

import (
	"WebFramework/orm/model"
	"database/sql"
)

type IValue interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(model *model.Model, entity any) IValue
