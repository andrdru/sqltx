package repo

import (
	"context"

	"github.com/andrdru/sqltx"
)

type (
	Address interface {
		sqltx.Tx

		TX(action func(txRepo Address) error) error

		CreateAddress(ctx context.Context, id int64, street string) (err error)
		GetAddressStreetByID(ctx context.Context, id int64) (name string, err error)
	}

	address struct {
		sqltx.Tx
	}
)

func NewAddress(db sqltx.QueryExecutor) Address {
	return &address{
		Tx: sqltx.NewTx(db),
	}
}

// TX simplify DoTransaction call for Address interface
func (m *address) TX(action func(txRepo Address) error) error {
	return m.DoTransaction(
		func(tx sqltx.QueryExecutor) sqltx.Tx { return NewAddress(tx) },
		func(txRepo sqltx.Tx) error { return action(txRepo.(Address)) },
	)
}

func (m *address) GetAddressStreetByID(ctx context.Context, id int64) (name string, err error) {
	const query = `SELECT street FROM addresses WHERE id = $1`

	err = m.DB().QueryRowContext(ctx, query, id).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (m *address) CreateAddress(ctx context.Context, id int64, street string) (err error) {
	const query = `INSERT INTO addresses(id, street) VALUES($1, $2)`

	_, err = m.DB().ExecContext(ctx, query, id, street)
	if err != nil {
		return err
	}

	return nil
}
