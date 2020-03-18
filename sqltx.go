package sqltx

import (
	"context"
	"database/sql"
	"errors"
)

var (
	// ErrDatabaseTypeInvalid throws when QueryExecutor not *sql.DB
	// are you try to run transaction in transaction?
	ErrDatabaseTypeInvalid = errors.New("wrong database struct type")
)

type (
	// QueryExecutor is execution interface of sql.Tx and sql.DB
	QueryExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}

	//Tx repo tx helper
	Tx interface {
		DoTx(fn func(Tx) error) (err error)
	}

	// tx tx helper implementation
	tx struct {
		db QueryExecutor
	}
)

// DoTx wrap repository calls into transaction
//
// func (m *myRepo) DoTransaction(action func(txRepo MyRepo) (err error)) (err error) {
//	return m.DoTx(func(tx sqltx.Tx) error {
//		var repo = NewMyRepo(m.db)
//		return action(repo)
//	})
// }
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

// withTx returns Tx instance with sql.Tx as QueryExecutor
func (r *tx) withTx(sqlTx *sql.Tx) Tx {
	return &tx{
		db: sqlTx,
	}
}
