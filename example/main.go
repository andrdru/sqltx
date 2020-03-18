package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/andrdru/sqltx/example/repo"
)

func main() {
	var err error
	var r = repo.NewMyRepo(&sql.DB{})

	err = r.DoTransaction(func(txRepo repo.MyRepo) (err error) {
		return txRepo.Ping(context.Background())
	})

	if err != nil {
		log.Fatal("tx failed:", err.Error())
	}
}
