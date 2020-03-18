package repo

import (
	"context"

	"github.com/andrdru/sqltx"
)

type (
	MyRepo interface {
		sqltx.Tx
		DoTransaction(action func(txRepo MyRepo) (err error)) (err error)
		Ping(ctx context.Context) (err error)
	}

	myRepo struct {
		sqltx.Tx
		db sqltx.QueryExecutor
	}
)

func NewMyRepo(db sqltx.QueryExecutor) MyRepo {
	return &myRepo{
		db: db,
	}
}

func (m *myRepo) DoTransaction(action func(txRepo MyRepo) (err error)) (err error) {
	return m.DoTx(func(tx sqltx.Tx) error {
		var repo = NewMyRepo(m.db)
		return action(repo)
	})
}

func (m *myRepo) Ping(ctx context.Context) (err error) {
	_, err = m.db.ExecContext(ctx, "SELECT 1")
	return err
}
