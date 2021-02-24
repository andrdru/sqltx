# Golang SQL transaction

Helper for sql transactions

## Install

`go get -u github.com/andrdru/sqltx`

## Usage

see [examples](/examples)

```go
// add helper to repo

User interface {
  sqltx.Tx

  TX(action func(txRepo User) error) error
}

// TX simplify DoTransaction call for User interface
func (m *user) TX(action func(txRepo User) error) error {
  return m.DoTransaction(
    func(tx sqltx.QueryExecutor) sqltx.Tx { return NewUser(tx) },
    func(txRepo sqltx.Tx) error { return action(txRepo.(User)) },
  )
}

```

```go
// call

err = userRepo.TX(func(txRepo repo.User) (err error) {
  err = txRepo.CreateUser(ctx, 1, "Vasya")
  if err != nil {
    return err
  }
  
  return nil
}

if err != nil {
  log.Fatalf("transaction rollback: %s", err)
}
```

