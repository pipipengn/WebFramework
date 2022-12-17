package orm

import (
	"WebFramework/orm/internal/valuer"
	"WebFramework/orm/model"
	"database/sql"
)

type DB struct {
	r          model.IRegistory
	sqldb      *sql.DB
	valCreator valuer.Creator
}

type DBOption func(db *DB)

func WithRegistory(r model.IRegistory) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func WithReflectCreator() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:          model.NewRegistory(),
		sqldb:      db,
		valCreator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func Open(driver, datasourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, datasourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

func MustOpen(driver, datasourceName string, opts ...DBOption) *DB {
	db, err := Open(driver, datasourceName, opts...)
	if err != nil {
		panic(err)
	}
	return db
}
