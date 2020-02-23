package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/andrdru/sqltx/repo"
)

func main() {
	var err error
	var r = repo.NewSomeRepo(&sql.DB{}, nil)

	err = r.DoTransaction(func(txRepo repo.SomeRepo) (err error) {
		return txRepo.Ping(context.Background())
	})

	if err != nil {
		log.Fatal("tx failed:", err.Error())
	}
}
