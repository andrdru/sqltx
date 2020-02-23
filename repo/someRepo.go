package repo

import (
	"context"
	"database/sql"

	"github.com/andrdru/sqltx/postgres"
)

type (
	SomeRepo interface {
		postgres.Tx
		DoTransaction(action func(txRepo SomeRepo) (err error)) (err error)
		Ping(ctx context.Context) (err error)
	}

	someRepo struct {
		postgres.Tx
		db *sql.DB
	}
)

func NewSomeRepo(db *sql.DB, tx postgres.Tx) SomeRepo {
	if tx == nil {
		tx = postgres.NewTx(db)
	}
	return &someRepo{
		db: db,
		Tx: tx,
	}
}

func (m *someRepo) DoTransaction(action func(txRepo SomeRepo) (err error)) (err error) {
	return m.DoTx(func(tx postgres.Tx) error {
		var repo = NewSomeRepo(m.db, tx)
		return action(repo)
	})
}

func (m *someRepo) Ping(ctx context.Context) (err error) {
	_, err = m.DBTx().ExecContext(ctx, "SELECT 1")
	return err
}
