package valuer

import (
	"WebFramework/orm/internal/errs"
	"WebFramework/orm/model"
	"database/sql"
	"reflect"
)

type reflectValue struct {
	model *model.Model
	t     any // the pointer of T
}

var _ Creator = NewReflectValue

func NewReflectValue(model *model.Model, t any) IValue {
	return &reflectValue{model: model, t: t}
}

func (r *reflectValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cols))
	valElems := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd, ok := r.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknowColumn(col)
		}
		v := reflect.New(fd.Type)
		vals = append(vals, v.Interface())
		valElems = append(valElems, v.Elem())
	}

	if err = rows.Scan(vals...); err != nil {
		return err
	}

	tval := reflect.ValueOf(r.t).Elem()
	for i, col := range cols {
		fd, ok := r.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknowColumn(col)
		}
		tval.FieldByName(fd.GoName).Set(valElems[i])
	}
	return nil
}
