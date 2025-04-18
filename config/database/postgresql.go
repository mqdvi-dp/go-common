package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config/database/dbc"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
)

var sqlMap = make(map[string]*sqlx.DB)

type SqlFuncOption func(*sqlOption)

type sqlOption struct {
	dsn               string
	maxIdleTime       time.Duration
	maxIdleConnection int
	maxConnection     int
}

func defaultSqlOption() sqlOption {
	return sqlOption{
		dsn:               env.GetString("DSN_MASTER", "postgres://127.0.0.1:5432"),
		maxIdleTime:       time.Duration(1) * time.Minute,
		maxIdleConnection: 5,
		maxConnection:     env.GetInt("DB_MAX_CONNECTION", 20),
	}
}

// sqlxInstance instance variables used by the sqlx database
type sqlxInstance struct {
	db dbc.SqlDbc
}

// Database get abstraction sqlx
func (s *sqlxInstance) Database() dbc.SqlDbc {
	return s.db
}

// Disconnect close the database connection
func (s *sqlxInstance) Disconnect(ctx context.Context) error {
	logger.RedBold("postgres: disconnecting...")
	defer fmt.Printf("\x1b[31;1mPostgres Disconnecting:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n")

	return s.db.Close()
}

// NewSqlxConnection creates a SqlxConnection
func NewSqlxConnection(opts ...SqlFuncOption) (abstract.SQLDatabase, error) {
	logger.YellowItalic("Load postgresql connection...")
	var (
		client *sqlx.DB
		err    error
	)

	// sql custom option
	opt := defaultSqlOption()
	// set option from parameters
	for _, o := range opts {
		o(&opt)
	}

	// if connection already declare, use that
	client, ok := sqlMap[opt.dsn]
	if ok {
		logger.GreenItalic("postgresql connected!")
		return &sqlxInstance{db: &dbc.DB{DB: client}}, nil
	}

	// initiate connection with postgres
	client, err = sqlx.Open("postgres", opt.dsn)
	if err != nil {
		panic(err)
	}

	// check connection
	err = client.Ping()
	if err != nil {
		panic(err)
	}

	client.SetMaxIdleConns(opt.maxIdleConnection)
	client.SetMaxOpenConns(opt.maxConnection)
	client.SetConnMaxIdleTime(opt.maxIdleTime)

	// store connection into hashMap
	sqlMap[opt.dsn] = client
	logger.GreenItalic("postgresql connected!")
	return &sqlxInstance{db: &dbc.DB{DB: client}}, nil
}

// SetSqlDSN sets the database dsn
func SetSqlDSN(dsn string) SqlFuncOption {
	return func(so *sqlOption) {
		so.dsn = dsn
	}
}

// SetSqlMaxIdleTime sets the max idle time for the connection
func SetSqlMaxIdleTime(maxIdleTime time.Duration) SqlFuncOption {
	return func(so *sqlOption) {
		so.maxIdleTime = maxIdleTime
	}
}

// SetSqlMaxIdleConnection sets the max idle connection for the connection
func SetSqlMaxIdleConnection(maxIdleConnection int) SqlFuncOption {
	return func(so *sqlOption) {
		so.maxIdleConnection = maxIdleConnection
	}
}

// SetSqlMaxConnection sets the max connection for the connection
func SetSqlMaxConnection(maxConnection int) SqlFuncOption {
	return func(so *sqlOption) {
		so.maxConnection = maxConnection
	}
}
