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

	// Tx repo tx helper
	Tx interface {
		DoTransaction(construct func(tx QueryExecutor) Tx, action func(txRepo Tx) error) (err error)
		DB() (db QueryExecutor)
	}

	// tx helper implementation
	tx struct {
		db QueryExecutor
	}
)

var (
	_ Tx            = &tx{}
	_ QueryExecutor = &sql.Tx{}
	_ QueryExecutor = &sql.DB{}
)

func NewTx(db QueryExecutor) Tx {
	return &tx{
		db: db,
	}
}

// DB get current QueryExecutor
func (r *tx) DB() (db QueryExecutor) {
	return r.db
}

// DoTransaction repository transactions helper
//
// // TX simplify DoTransaction call for MyRepo interface
// func (m *myRepo) TX(action func(txRepo MyRepo) error) error {
//	return m.DoTransaction(
//		func(tx sqltx.QueryExecutor) sqltx.Tx { return NewMyRepo(tx) },
//		func(txRepo sqltx.Tx) error { return action(txRepo.(MyRepo)) },
//	)
// }
//
func (r *tx) DoTransaction(construct func(tx QueryExecutor) Tx, action func(txRepo Tx) error) (err error) {
	return r.doTx(func(tx Tx) error {
		return action(construct(tx.DB()))
	})
}

// DoTx wrap repository calls into transaction
func (r *tx) doTx(fn func(Tx) error) (err error) {
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
