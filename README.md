# Golang SQL transaction

Helper for sql transactions

## Install

`go get -u github.com/andrdru/sqltx`

## Usage

```go
// add helper to repo

User interface {
  sqltx.Tx
  DoTransaction(action func (txRepo User) (err error)) (err error)
}

// DoTransaction repository transactions helper
func (m *user) DoTransaction(action func(txRepo User) (err error)) (err error) {
  return m.DoTx(func(tx sqltx.Tx) error {
    var repo = NewUser(tx.DB())
    return action(repo)
  })
}
```

```go
// call

err = userRepo.DoTransaction(func(txRepo repo.User) (err error) {
  err = txRepo.CreateUser(ctx, 1, "Vasya")
  if err != nil {
    log.Fatalf("could not get user: %s", name)
  }
  
  return nil
}

if err != nil {
  log.Fatalf("transaction: %s", err)
}
```

full example in [examples](/examples)