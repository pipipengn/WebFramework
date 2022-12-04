package orm

import (
	"WebFramework/orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			assert.Equalf(t, tt.want, formatName(tt.s), "formatName(%v)", tt.s)
		})
	}
}

func Test_parseModel(t *testing.T) {

	tests := []struct {
		name    string
		entity  any
		want    *model
		wantErr error
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			want: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Id":        {colName: "id"},
					"Age":       {colName: "age"},
					"FirstName": {colName: "first_name"},
					"LastName":  {colName: "last_name"},
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseModel(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want.tableName, res.tableName)
			assert.Equal(t, tt.want.fields, res.fields)
		})
	}
}
