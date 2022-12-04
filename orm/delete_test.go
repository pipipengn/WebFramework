package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   QueryBuilder
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no where",
			builder: NewDeleter[TestModel]().From("test_model"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: NewDeleter[TestModel]().Where(Col("Id").Eq(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id`=?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: NewDeleter[TestModel]().From("test_model").Where(Col("Id").Eq(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `id`=?;",
				Args: []any{16},
			},
		},
		{
			name:    "and",
			builder: NewDeleter[TestModel]().Where(Col("Id").Eq(16).And(Col("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`id`=?)AND(`age`=?);",
				Args: []any{16, 18},
			},
		},
		{
			name:    "or",
			builder: NewDeleter[TestModel]().Where(Col("Id").Eq(16).Or(Col("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`id`=?)OR(`age`=?);",
				Args: []any{16, 18},
			},
		},
		{
			name:    "not",
			builder: NewDeleter[TestModel]().Where(Not(Col("Id").Eq(16))),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE NOT(`id`=?);",
				Args: []any{16},
			},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
