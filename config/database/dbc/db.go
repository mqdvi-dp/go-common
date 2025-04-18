package dbc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (d *DB) Queryx(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:Queryx")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	rows, err = d.DB.QueryxContext(ctx, query, args...)
	return
}

func (d *DB) QueryRowx(ctx context.Context, query string, args ...interface{}) (row *sqlx.Row) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:QueryRowx")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	row = d.DB.QueryRowxContext(ctx, query, args...)
	return
}

func (d *DB) NamedExec(ctx context.Context, query string, args interface{}) (row sql.Result, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:NamedExec")
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

func (d *DB) Exec(ctx context.Context, query string, args ...interface{}) (row sql.Result, err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:Exec")
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

func (d *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:Get")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	err = d.DB.GetContext(ctx, dest, query, args...)
	return
}

func (d *DB) GetWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:GetWithIn")
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

func (d *DB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:Select")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

	// log tracer
	trace.Log("query", query)
	trace.Log("arguments", args)

	err = d.DB.SelectContext(ctx, dest, query, args...)
	return
}

func (d *DB) SelectWithIn(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:SelectWithIn")
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

func (d *DB) Preparex(ctx context.Context, query string, args ...interface{}) (err error) {
	var log logger.Database
	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:Preparex")
	defer func() {
		log.Store(ctx)
		trace.SetError(err)
		trace.Finish()
	}()
	log = logger.DB(logger.Sql, query, args...)

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

	// execute the query
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		trace.SetError(err)

		return
	}

	// get rows affected after execute the query
	rowsAffected, _ := result.RowsAffected()
	// when rows affected is more than 0, it means, execute the query is success
	if rowsAffected > 0 {
		return nil
	}

	// otherwise, send the error
	err = ErrNoRowsAffected
	trace.SetError(err)
	return
}

func (d *DB) Rebind(query string) string {
	return d.DB.Rebind(query)
}

func (d *DB) MustBegin() *Tx {
	tx := d.DB.MustBegin()

	return &Tx{tx}
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func (d *DB) Ping() error {
	return d.DB.Ping()
}

func (d *DB) StartTransaction(ctx context.Context, txFunc func(context.Context, SqlDbc) error) error {
	var err error
	var tx *Tx

	trace, ctx := tracer.StartTraceWithContext(ctx, "Sql:StartTransaction")
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
			trace.SetError(err)

			if tx != nil {
				_ = tx.Rollback()
			}

			return
		} else if err != nil {
			_ = tx.Rollback()
			trace.SetError(err)
			return
		}

		_ = tx.Commit()
		trace.Finish()
	}()

	tx = d.MustBegin()

	err = txFunc(ctx, tx)
	return err
}

func (d *DB) Rollback() error {
	return nil
}

func (d *DB) Commit() error {
	return nil
}
