package orm

import (
	"WebFramework/orm/internal/errs"
	"WebFramework/orm/model"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	db := &DB{
		r: model.NewRegistory(),
	}
	tests := []struct {
		name      string
		builder   QueryBuilder
		wantQuery *Query
		wantError error
	}{
		{
			name:    "no from",
			builder: NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](db).From("test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_model;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "empty from",
			builder: NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "from db",
			builder: NewSelector[TestModel](db).From("test_db.test_model"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM test_db.test_model;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](db).Where(Col("Age").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `age`=?;",
				Args: []any{18},
			},
			wantError: nil,
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](db).Where(Not(Col("Age").Eq(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE NOT(`age`=?);",
				Args: []any{18},
			},
			wantError: nil,
		},
		{
			name:    "and",
			builder: NewSelector[TestModel](db).Where(Col("Age").Eq(18).And(Col("LastName").Eq("ppp"))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age`=?)AND(`last_name`=?);",
				Args: []any{18, "ppp"},
			},
			wantError: nil,
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](db).Where(),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
			wantError: nil,
		},
		{
			name:      "invalid column",
			builder:   NewSelector[TestModel](db).Where(Col("xxx").Eq(123)),
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

func TestSelector_Get(t *testing.T) {
	sqldb, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(sqldb)
	require.NoError(t, err)

	// query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))
	// no data
	rows := sqlmock.NewRows([]string{"id", "first_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
	// select
	rows = sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
	rows.AddRow("1", "18", "ppp", "123")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	tests := []struct {
		name    string
		s       *Selector[TestModel]
		want    *TestModel
		wantErr error
	}{
		{
			name:    "invalid query",
			s:       NewSelector[TestModel](db).Where(Col("xxx").Eq(123)),
			wantErr: errs.NewErrUnknowField("xxx"),
		},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(Col("Id").Eq(123)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no data",
			s:       NewSelector[TestModel](db).Where(Col("Id").Eq(123)),
			wantErr: ErrNoRows,
		},
		{
			name: "select",
			s:    NewSelector[TestModel](db).Where(Col("Id").Eq(1)),
			want: &TestModel{
				Id:        1,
				Age:       18,
				FirstName: "ppp",
				LastName: &sql.NullString{
					String: "123",
					Valid:  true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.s.Get(context.Background())
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestSelector_Select(t *testing.T) {
	mockdb, _, _ := sqlmock.New()
	db, _ := OpenDB(mockdb)

	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "1",
			s:    NewSelector[TestModel](db).Select(Col("Id"), Col("Age")),
			want: &Query{
				SQL: "SELECT `id`,`age` FROM `test_model`;",
			},
		},
		{
			name:    "invalid",
			s:       NewSelector[TestModel](db).Select(Col("Iddd")),
			wantErr: errs.NewErrUnknowField("Iddd"),
		},
		{
			name: "avg",
			s:    NewSelector[TestModel](db).Select(Avg("Age")),
			want: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "sum",
			s:    NewSelector[TestModel](db).Select(Sum("Age")),
			want: &Query{
				SQL: "SELECT SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "max",
			s:    NewSelector[TestModel](db).Select(Max("Age")),
			want: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "min",
			s:    NewSelector[TestModel](db).Select(Min("Age")),
			want: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},
		{
			name: "count",
			s:    NewSelector[TestModel](db).Select(Count("Age")),
			want: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},
		{
			name:    "aggregate invalid",
			s:       NewSelector[TestModel](db).Select(Count("Ageddddd")),
			wantErr: errs.NewErrUnknowField("Ageddddd"),
		},
		{
			name: "aggregate many",
			s:    NewSelector[TestModel](db).Select(Count("Age"), Avg("Age")),
			want: &Query{
				SQL: "SELECT COUNT(`age`),AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw",
			s:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT first_name)")),
			want: &Query{
				SQL: "SELECT COUNT(DISTINCT first_name) FROM `test_model`;",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, query)
		})
	}
}

func TestRawExpr(t *testing.T) {
	mockdb, _, _ := sqlmock.New()
	db, _ := OpenDB(mockdb)

	tests := []struct {
		name    string
		s       QueryBuilder
		want    *Query
		wantErr error
	}{
		{
			name: "raw1",
			s:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT first_name)")),
			want: &Query{
				SQL: "SELECT COUNT(DISTINCT first_name) FROM `test_model`;",
			},
		},
		{
			name: "raw2",
			s:    NewSelector[TestModel](db).Where(Raw("id < ?", 12).AsPredicate()),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE id < ?;",
				Args: []any{12},
			},
		},
		{
			name: "raw3",
			s:    NewSelector[TestModel](db).Where(Col("Id").Eq(Raw("age + ?", 1))),
			want: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id`=age + ?;",
				Args: []any{1},
			},
		},
		{
			name: "raw4",
			s:    NewSelector[TestModel](db).Where(Col("Id").Eq(Raw("age + 1"))),
			want: &Query{
				SQL: "SELECT * FROM `test_model` WHERE `id`=age + 1;",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := tt.s.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, query)
		})
	}
}
