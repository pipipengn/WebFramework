package model

import (
	"WebFramework/orm/internal/errs"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type TestModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}

func Test_formatName(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "1",
			s:    "TestModel",
			want: "test_model",
		},
		{
			name: "2",
			s:    "testModel",
			want: "test_model",
		},
		{
			name: "3",
			s:    "test_model",
			want: "test_model",
		},
		{
			name: "4",
			s:    "testmodel",
			want: "testmodel",
		},
		{
			name: "5",
			s:    "Test1Model",
			want: "test1_model",
		},
		{
			name: "6",
			s:    "ID",
			want: "i_d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, underscoreName(tt.s), "underscoreName(%v)", tt.s)
		})
	}
}

func Test_Register(t *testing.T) {

	tests := []struct {
		name    string
		entity  any
		want    *Model
		wantErr error
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			want: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						GoName:  "Id",
						ColName: "id",
						Type:    reflect.TypeOf(TestModel{}).Field(0).Type,
					},
					"Age": {
						GoName:  "Age",
						ColName: "age",
						Type:    reflect.TypeOf(TestModel{}).Field(1).Type,
					},
					"FirstName": {
						GoName:  "FirstName",
						ColName: "first_name",
						Type:    reflect.TypeOf(TestModel{}).Field(2).Type,
					},
					"LastName": {
						GoName:  "LastName",
						ColName: "last_name",
						Type:    reflect.TypeOf(TestModel{}).Field(3).Type,
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "struct",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "slice",
			entity:  []TestModel{},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "primitive",
			entity:  0,
			wantErr: errs.ErrPointerOnly,
		},
	}

	r := NewRegistory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := r.Register(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want.TableName, res.TableName)
			for _, v := range res.FieldMap {
				v.Offset = 0
			}
			assert.Equal(t, tt.want.FieldMap, res.FieldMap)
		})
	}
}

func Test_Get(t *testing.T) {

	type TagTable1 struct {
		FirstName string `orm:"column=first_name_t"`
	}

	type TagTableEmptyColumn struct {
		FirstName string `orm:"column="`
	}

	type TagTableIgnore struct {
		FirstName string `orm:"xxx=123"`
	}

	tests := []struct {
		name    string
		entity  any
		want    *Model
		wantErr error
	}{
		{
			name:   "1",
			entity: &TagTable1{},
			want: &Model{
				TableName: "tag_table1",
				FieldMap: map[string]*Field{
					"FirstName": {
						GoName:  "FirstName",
						ColName: "first_name_t",
						Type:    reflect.TypeOf(TagTable1{}).Field(0).Type,
					},
				},
			},
		},
		{
			name:   "empty column",
			entity: &TagTableEmptyColumn{},
			want: &Model{
				TableName: "tag_table_empty_column",
				FieldMap: map[string]*Field{
					"FirstName": {
						GoName:  "FirstName",
						ColName: "first_name",
						Type:    reflect.TypeOf(TagTableEmptyColumn{}).Field(0).Type,
					},
				},
			},
		},
		{
			name: "invalid tag content",
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			name:   "tag ignore",
			entity: &TagTableIgnore{},
			want: &Model{
				TableName: "tag_table_ignore",
				FieldMap: map[string]*Field{
					"FirstName": {
						GoName:  "FirstName",
						ColName: "first_name",
						Type:    reflect.TypeOf(TagTableIgnore{}).Field(0).Type,
					},
				},
			},
		},
		{
			name:   "custom table name",
			entity: &CustomtableName{},
			want: &Model{
				TableName: "CustomtableName_t",
				FieldMap:  map[string]*Field{},
			},
		},
		{
			name:   "custom table name ptr",
			entity: &CustomtableNamePtr{},
			want: &Model{
				TableName: "CustomtableName_t",
				FieldMap:  map[string]*Field{},
			},
		},
		{
			name:   "custom table name empty string",
			entity: &CustomTableNameEmptyString{},
			want: &Model{
				TableName: "custom_table_name_empty_string",
				FieldMap:  map[string]*Field{},
			},
		},
	}

	r := NewRegistory()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := r.Get(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want.TableName, res.TableName)
			assert.Equal(t, tt.want.FieldMap, res.FieldMap)
			// 测试有没有缓存住
			typ := reflect.TypeOf(tt.entity)
			cache, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, cache, res)
		})
	}
}

type CustomtableName struct{}

func (c CustomtableName) TableName() string {
	return "CustomtableName_t"
}

type CustomtableNamePtr struct{}

func (c *CustomtableNamePtr) TableName() string {
	return "CustomtableName_t"
}

type CustomTableNameEmptyString struct{}

func (c *CustomTableNameEmptyString) TableName() string {
	return ""
}

func TestWithTableName(t *testing.T) {
	r := NewRegistory()
	m, err := r.Register(&OptionTable{}, WithTableName("OptionTable1"))
	require.NoError(t, err)
	assert.Equal(t, "OptionTable1", m.TableName)
}

func TestWithColumnName(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		opts    []Option
		want    *Model
		wantErr error
	}{
		{
			name:   "1",
			entity: &OptionTable{},
			opts:   []Option{WithColunmName("FirstName", "FirstName1")},
			want: &Model{
				TableName: "option_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						GoName:  "FirstName",
						ColName: "FirstName1",
						Type:    reflect.TypeOf(OptionTable{}).Field(0).Type,
					},
				},
				ColumnMap: map[string]*Field{
					"FirstName1": {
						GoName:  "FirstName",
						ColName: "FirstName1",
						Type:    reflect.TypeOf(OptionTable{}).Field(0).Type,
					},
				},
			},
		},
		{
			name:    "2",
			entity:  &OptionTable{},
			opts:    []Option{WithColunmName("FirstName1", "FirstName1")},
			wantErr: errs.NewErrUnknowField("FirstName1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistory()
			m, err := r.Register(tt.entity, tt.opts...)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, m)
		})
	}
}

type OptionTable struct {
	FirstName string
}
