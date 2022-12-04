package orm

import (
	"WebFramework/orm/internal/errs"
	"reflect"
	"strings"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

// entity只能使用一级指针
func parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()

	numField := typ.NumField()
	model := &model{
		tableName: formatName(typ.Name()),
		fields:    make(map[string]*field, numField),
	}

	for i := 0; i < numField; i++ {
		fieldType := typ.Field(i)
		model.fields[fieldType.Name] = &field{
			colName: formatName(fieldType.Name),
		}
	}
	return model, nil
}

func formatName(s string) string {
	b := strings.Builder{}
	for i, r := range s {
		if !unicode.IsUpper(r) {
			b.WriteRune(r)
			continue
		}
		if i != 0 {
			b.WriteRune('_')
		}
		b.WriteString(strings.ToLower(string(r)))
	}
	return b.String()
}
