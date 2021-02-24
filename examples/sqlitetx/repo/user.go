package repo

import (
	"context"

	"github.com/andrdru/sqltx"
)

type (
	User interface {
		sqltx.Tx
		DoTransaction(action func(txRepo User) (err error)) (err error)

		CreateUser(ctx context.Context, id int64, name string) (err error)
		GetUserNameByID(ctx context.Context, id int64) (name string, err error)
	}

	user struct {
		sqltx.Tx
	}
)

func NewUser(db sqltx.QueryExecutor) User {
	return &user{
		Tx: sqltx.NewTx(db),
	}
}

// DoTransaction repository transactions helper
func (m *user) DoTransaction(action func(txRepo User) (err error)) (err error) {
	return m.DoTx(func(tx sqltx.Tx) error {
		var repo = NewUser(tx.DB())
		return action(repo)
	})
}

func (m *user) GetUserNameByID(ctx context.Context, id int64) (name string, err error) {
	const query = `SELECT name FROM users WHERE id = ?`

	err = m.DB().QueryRowContext(ctx, query, id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (m *user) CreateUser(ctx context.Context, id int64, name string) (err error) {
	const query = `INSERT INTO users(id, name) VALUES(?, ?)`

	_, err = m.DB().ExecContext(ctx, query, id, name)
	if err != nil {
		return err
	}

	return nil
}
