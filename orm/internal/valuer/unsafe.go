package valuer

import (
	"WebFramework/orm/internal/errs"
	"WebFramework/orm/model"
	"database/sql"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	t     any // the pointer of T
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(model *model.Model, t any) IValue {
	return &unsafeValue{model: model, t: t}
}

func (u *unsafeValue) SetColumns(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cols))
	addr := reflect.ValueOf(u.t).UnsafePointer()
	for _, col := range cols {
		fd, ok := u.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknowColumn(col)
		}
		val := reflect.NewAt(fd.Type, unsafe.Pointer(uintptr(addr)+fd.Offset))
		vals = append(vals, val.Interface())
	}

	err = rows.Scan(vals...)
	return err
}
