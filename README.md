# Golang SQL transaction

Helper for sql transactions

## Install
`go get -u github.com/andrdru/sqltx`

## Usage

```go
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

func (m *myRepo) DoTransaction(action func(txRepo MyRepo) (err error)) (err error) {
	return m.DoTx(func(tx sqltx.Tx) error {
		var repo = NewMyRepo(m.db)
		return action(repo)
	})
}
```

```go
var err error
var r = repo.NewMyRepo(&sql.DB{})

err = r.DoTransaction(func(txRepo repo.MyRepo) (err error) {
	return txRepo.Ping(context.Background())
})

if err != nil {
	log.Fatal("tx failed:", err.Error())
}
```