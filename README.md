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
	}
)

func NewMyRepo(db sqltx.QueryExecutor) MyRepo {
	return &myRepo{
		Tx: sqltx.NewTx(db),
	}
}

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