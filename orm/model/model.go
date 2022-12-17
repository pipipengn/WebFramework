package model

import (
	"WebFramework/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagOrmKey    = "orm"
	tagColumnKey = "column"
)

type Model struct {
	TableName string
	// GoName -> Field
	FieldMap map[string]*Field
	// ColName -> Field
	ColumnMap map[string]*Field
}

type Field struct {
	GoName  string
	ColName string
	Type    reflect.Type
	Offset  uintptr
}

type Option func(model *Model) error

type IRegistory interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...Option) (*Model, error)
}

// Registory 缓存元数据 每一个orm的struct都会对应一个model
type Registory struct {
	models sync.Map // map[reflect.Type]*Model
}

func NewRegistory() *Registory {
	return &Registory{}
}

func (r *Registory) Get(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}

	m, err := r.Register(entity)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), err
}

// Register entity只能使用一级指针
func (r *Registory) Register(entity any, opts ...Option) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()

	var tableName string
	if tableNameInterface, ok := entity.(TableName); ok {
		tableName = tableNameInterface.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}

	model := &Model{
		TableName: tableName,
		FieldMap:  make(map[string]*Field, numField),
		ColumnMap: make(map[string]*Field, numField),
	}

	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		tags, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags[tagColumnKey]
		if colName == "" {
			colName = underscoreName(fd.Name)
		}
		f := &Field{
			GoName:  fd.Name,
			ColName: colName,
			Type:    fd.Type,
			Offset:  fd.Offset,
		}
		model.FieldMap[fd.Name] = f
		model.ColumnMap[colName] = f
	}

	for _, opt := range opts {
		err := opt(model)
		if err != nil {
			return nil, err
		}
	}
	return model, nil
}

func WithTableName(name string) Option {
	return func(m *Model) error {
		m.TableName = name
		return nil
	}
}

func WithColunmName(goName, colName string) Option {
	return func(m *Model) error {
		fd, ok := m.FieldMap[goName]
		if !ok {
			return errs.NewErrUnknowField(goName)
		}
		oldColName := fd.ColName
		fd.ColName = colName

		fd = m.ColumnMap[oldColName]
		m.ColumnMap[colName] = fd
		delete(m.ColumnMap, oldColName)
		return nil
	}
}

func (r *Registory) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup(tagOrmKey)
	if !ok {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		pair = strings.Trim(pair, " ")
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return map[string]string{}, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

func underscoreName(s string) string {
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

// TableName 自定义表名
type TableName interface {
	TableName() string
}
