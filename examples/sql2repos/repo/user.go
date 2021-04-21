package repo

import (
	"context"

	"github.com/andrdru/sqltx"
)

type (
	User interface {
		sqltx.Tx

		TX(action func(txRepo User) error) error

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

// TX simplify DoTransaction call for User interface
func (m *user) TX(action func(txRepo User) error) error {
	return m.DoTransaction(
		func(tx sqltx.QueryExecutor) sqltx.Tx { return NewUser(tx) },
		func(txRepo sqltx.Tx) error { return action(txRepo.(User)) },
	)
}

func (m *user) GetUserNameByID(ctx context.Context, id int64) (name string, err error) {
	const query = `SELECT name FROM users WHERE id = $1`

	err = m.DB().QueryRowContext(ctx, query, id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (m *user) CreateUser(ctx context.Context, id int64, name string) (err error) {
	const query = `INSERT INTO users(id, name) VALUES($1, $2)`

	_, err = m.DB().ExecContext(ctx, query, id, name)
	if err != nil {
		return err
	}

	return nil
}
