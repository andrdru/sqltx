package postgres

import (
	"context"
	"database/sql"
)

type (
	// sql.Tx interface
	TxExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	// Repo tx helper
	Tx interface {
		DoTx(fn func(Tx) error) (err error)
		WithTx(sqlTx *sql.Tx) Tx
		DBTx() TxExecutor
	}

	tx struct {
		db *sql.DB
		tx *sql.Tx
	}
)

func NewTx(db *sql.DB) Tx {
	return &tx{
		db: db,
	}
}

func (r *tx) DBTx() TxExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *tx) DoTx(fn func(Tx) error) (err error) {
	var tx *sql.Tx
	if tx, err = r.db.Begin(); err != nil {
		return err
	}

	if err = fn(r.WithTx(tx)); err != nil {
		var errRollback error
		if errRollback = tx.Rollback(); errRollback != nil {
			return errRollback
		}

		return err
	}

	return tx.Commit()
}

func (r *tx) WithTx(sqlTx *sql.Tx) Tx {
	return &tx{
		db: r.db,
		tx: sqlTx,
	}
}
