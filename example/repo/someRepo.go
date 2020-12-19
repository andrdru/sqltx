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
	}
)

func NewMyRepo(db sqltx.QueryExecutor) MyRepo {
	return &myRepo{
		Tx: sqltx.NewTx(db),
	}
}

// DoTransaction repository transactions helper
func (m *myRepo) DoTransaction(action func(txRepo MyRepo) (err error)) (err error) {
	return m.DoTx(func(tx sqltx.Tx) error {
		var repo = NewMyRepo(m.DB())
		return action(repo)
	})
}

func (m *myRepo) Ping(ctx context.Context) (err error) {
	_, err = m.DB().ExecContext(ctx, "SELECT 1")
	return err
}
