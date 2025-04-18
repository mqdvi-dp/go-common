package dbc

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

// GormDB is an instance with field struct non-transaction connection with gorm.io v2
type GormDB struct {
	DB *gorm.DB
}

// GormTx is an instance with field struct transaction connection with gorm.io v2
type GormTx struct {
	DB *gorm.Tx
}

// DB is an instance with field struct non-transactional connection
type DB struct {
	DB *sqlx.DB
}

// Tx is istance with field struct transactional connection
type Tx struct {
	DB *sqlx.Tx
}

type SqlDbc interface {
	// Queryx returns a lot of data
	Queryx(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryRowsx return a single row data
	QueryRowx(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// Preparex to prepare statement query and execute the query with parameters
	Preparex(ctx context.Context, query string, args ...interface{}) error

	// Get return a single row data
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// GetWithIn return a single row data by conditional IN
	GetWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Select returns a lot of data
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// SelectWithIn retruns a lot data with query `in`
	SelectWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Rebind query sqlx.Name into paramater
	Rebind(query string) string

	// Insert data exec
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Insert data with named exec
	NamedExec(ctx context.Context, query string, args interface{}) (sql.Result, error)

	// MustBegin create transactional connection
	MustBegin() *Tx

	// Close connection
	Close() error

	// Check connection
	Ping() error

	// StartTransaction starts a transaction query
	StartTransaction(ctx context.Context, txFunc func(context.Context, SqlDbc) error) error

	// Transaction if you want support transactional connection
	Transaction
}

// Transaction implements all method for transactions
type Transaction interface {
	// Rollback a transaction
	Rollback() error

	// Commit a transaction
	Commit() error
}
