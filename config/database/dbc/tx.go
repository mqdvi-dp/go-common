package dbc

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (d *Tx) Queryx(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:Queryx")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.SqlTx, query, args...)

	rows, err = d.DB.QueryxContext(ctx, query, args...)
	return
}

func (d *Tx) QueryRowx(ctx context.Context, query string, args ...interface{}) (row *sqlx.Row) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:QueryRowx")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()
	log = logger.DB(logger.SqlTx, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	row = d.DB.QueryRowxContext(ctx, query, args...)
	return
}

func (d *Tx) Exec(ctx context.Context, query string, args ...interface{}) (row sql.Result, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:Exec")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("exec", query)
	trace.Log("arguments", args)
	res, err := d.DB.ExecContext(ctx, query, args...)
	if err != nil {
		trace.SetError(err)
		return nil, err
	}
	return res, err
}

func (d *Tx) NamedExec(ctx context.Context, query string, args interface{}) (row sql.Result, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:NamedExec")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args)

	// log tracer
	trace.Log("exec", query)
	trace.Log("arguments", args)
	res, err := d.DB.NamedExecContext(ctx, query, args)
	if err != nil {
		trace.SetError(err)
		return nil, err
	}
	return res, err
}

func (d *Tx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:Get")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.SqlTx, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	err = d.DB.GetContext(ctx, dest, query, args...)
	return
}

func (d *Tx) GetWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:GetWithIn")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	// create parameters "in"
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return
	}

	// rebind the query
	query = d.Rebind(query)

	// exec the query
	err = d.DB.GetContext(ctx, dest, query, args...)
	return
}

func (d *Tx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:Select")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.SqlTx, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	err = d.DB.SelectContext(ctx, dest, query, args...)
	return
}

func (d *Tx) SelectWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:SelectWithIn")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	// create parameters "in"
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return
	}

	// rebind the query
	query = d.Rebind(query)

	// exec the query
	err = d.DB.SelectContext(ctx, dest, query, args...)
	return
}

func (d *Tx) Preparex(ctx context.Context, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "SqlTx:Preparex")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.SqlTx, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	// create statement query
	stmt, err := d.DB.PreparexContext(ctx, query)
	if err != nil {
		trace.SetError(err)

		return
	}
	defer stmt.Close()

	// execute the query with parameters
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		trace.SetError(err)

		return
	}

	// get rows affected
	rowsAffected, _ := result.RowsAffected()
	// if rows affected more than 0, it means, success exec the query
	if rowsAffected > 0 {
		return
	}

	// otherwise, will send error
	err = ErrNoRowsAffected
	trace.SetError(err)
	return
}

func (d *Tx) Rebind(query string) string {
	return d.DB.Rebind(query)
}

func (d *Tx) MustBegin() *Tx {
	return d
}

func (d *Tx) Close() error {
	return nil
}

func (d *Tx) Ping() error {
	return nil
}

func (d *Tx) StartTransaction(ctx context.Context, txFunc func(context.Context, SqlDbc) error) error {
	return nil
}

func (d *Tx) Rollback() error {
	return d.DB.Rollback()
}

func (d *Tx) Commit() error {
	return d.DB.Commit()
}
