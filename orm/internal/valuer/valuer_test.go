package valuer

import (
	"WebFramework/orm/model"
	"database/sql"
	"database/sql/driver"
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

func TestSetColumns(t *testing.T) {
	test := func(t *testing.T, creator Creator) {
		tests := []struct {
			name    string
			entity  any
			rows    *sqlmock.Rows
			want    any
			wantErr error
		}{
			{
				name:   "set column",
				entity: &TestModel{},
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
					rows.AddRow("1", 18, "Tom", "Jerry")
					return rows
				}(),
				want: &TestModel{
					Id:        1,
					Age:       18,
					FirstName: "Tom",
					LastName: &sql.NullString{
						String: "Jerry",
						Valid:  true,
					},
				},
			},
			{
				name:   "set column no order",
				entity: &TestModel{},
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
					rows.AddRow("1", "Tom", "Jerry", 18)
					return rows
				}(),
				want: &TestModel{
					Id:        1,
					Age:       18,
					FirstName: "Tom",
					LastName: &sql.NullString{
						String: "Jerry",
						Valid:  true,
					},
				},
			},
			{
				name:   "partial column",
				entity: &TestModel{},
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "last_name"})
					rows.AddRow("1", "Jerry")
					return rows
				}(),
				want: &TestModel{
					Id:       1,
					LastName: &sql.NullString{String: "Jerry", Valid: true},
				},
			},
		}

		r := model.NewRegistory()
		mockDB, mock, _ := sqlmock.New()
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockRows := tt.rows
				mock.ExpectQuery("select .*").WillReturnRows(mockRows)
				rows, err := mockDB.Query("select .*")
				rows.Next()

				m, err := r.Get(tt.entity)
				require.NoError(t, err)
				value := creator(m, tt.entity)
				require.NoError(t, err)
				err = value.SetColumns(rows)
				assert.Equal(t, tt.wantErr, err)
				if err != nil {
					return
				}
				assert.Equal(t, tt.want, tt.entity)
			})
		}
	}

	test(t, NewReflectValue)
	test(t, NewUnsafeValue)
}

func BenchmarkSetColumns(b *testing.B) {
	r := model.NewRegistory()
	mockDB, mock, _ := sqlmock.New()

	benchmark := func(b *testing.B, creator Creator) {
		mockRows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
		vals := []driver.Value{"1", "Tom", "Jerry", 18}
		for i := 0; i < b.N; i++ {
			mockRows.AddRow(vals...)
		}
		mock.ExpectQuery("select .*").WillReturnRows(mockRows)
		rows, err := mockDB.Query("select .*")
		require.NoError(b, err)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows.Next()
			m, err := r.Get(&TestModel{})
			require.NoError(b, err)
			value := creator(m, &TestModel{})
			require.NoError(b, err)
			_ = value.SetColumns(rows)
		}
	}

	b.Run("reflect", func(b *testing.B) {
		benchmark(b, NewReflectValue)
	})

	b.Run("unsafe", func(b *testing.B) {
		benchmark(b, NewReflectValue)
	})
}
