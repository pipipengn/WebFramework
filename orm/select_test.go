package orm

import (
	"WebFramework/orm/internal/errs"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	tests := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantError error
	}{
		{
			name:    "no from",
			builder: NewSelector[TestModel](),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "from",
			builder: NewSelector[TestModel]().From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_model;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel]().From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "from db",
			builder: NewSelector[TestModel]().From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_db.test_model;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "where",
			builder: NewSelector[TestModel]().Where(Col("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age`=?;",
				Args: []any{18},
			},
			wantError: nil,
		},
		{
			name:    "not",
			builder: NewSelector[TestModel]().Where(Not(Col("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE NOT(`age`=?);",
				Args: []any{18},
			},
			wantError: nil,
		},
		{
			name:    "and",
			builder: NewSelector[TestModel]().Where(Col("Age").Eq(18).And(Col("LastName").Eq("ppp"))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age`=?)AND(`last_name`=?);",
				Args: []any{18, "ppp"},
			},
			wantError: nil,
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel]().Where(),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:      "invalid column",
			builder:   NewSelector[TestModel]().Where(Col("xxx").Eq(123)),
			wantError: errs.NewErrUnknowField("xxx"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.builder.Build()
			assert.Equal(t, tt.wantError, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantQuery, query)
		})
	}
}

type TestModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}
