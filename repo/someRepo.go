package repo

import (
	"context"

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
		db postgres.QueryExecutor
	}
)

func NewSomeRepo(db postgres.QueryExecutor) SomeRepo {
	return &someRepo{
		db: db,
	}
}

func (m *someRepo) DoTransaction(action func(txRepo SomeRepo) (err error)) (err error) {
	return m.DoTx(func(tx postgres.Tx) error {
		var repo = NewSomeRepo(m.db)
		return action(repo)
	})
}

func (m *someRepo) Ping(ctx context.Context) (err error) {
	_, err = m.db.ExecContext(ctx, "SELECT 1")
	return err
}
