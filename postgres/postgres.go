package postgres

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// are you try to run transaction in transaction ?
	ErrDatabaseTypeInvalid = errors.New("wrong database struct type")
)

type (
	// sql.Tx interface
	QueryExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// Repo tx helper
	Tx interface {
		DoTx(fn func(Tx) error) (err error)
	}

	tx struct {
		db QueryExecutor
	}
)

func (r *tx) DoTx(fn func(Tx) error) (err error) {
	var db *sql.DB
	var ok bool

	if db, ok = r.db.(*sql.DB); !ok {
		return ErrDatabaseTypeInvalid
	}

	var tx *sql.Tx
	if tx, err = db.Begin(); err != nil {
		return err
	}

	if err = fn(r.withTx(tx)); err != nil {
		var errRollback error
		if errRollback = tx.Rollback(); errRollback != nil {
			return errRollback
		}

		return err
	}

	return tx.Commit()
}

func (r *tx) withTx(sqlTx *sql.Tx) Tx {
	return &tx{
		db: sqlTx,
	}
}
