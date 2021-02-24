package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library

	"github.com/andrdru/sqltx/examples/sqlitetx/repo"
)

const dbname = "sqlite.db"

func main() {
	var ctx = context.Background()

	_ = os.Remove(dbname) // drop old database on restart

	var db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatalf("could not open db: %s\n", err)
	}

	defer func() {
		_ = db.Close()
	}()

	createTable(ctx, db)

	var userRepo = repo.NewUser(db)

	fmt.Printf("example commit:\n")
	commitTX(ctx, userRepo)

	fmt.Printf("\n\n")

	fmt.Printf("example rollback:\n")
	rollbackTX(ctx, userRepo)
}

func createTable(ctx context.Context, db *sql.DB) {
	const query = `CREATE TABLE users(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT);`

	var _, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("could not create table: %s\n", err)
	}
}

func commitTX(ctx context.Context, userRepo repo.User) {
	var name string
	var err = userRepo.DoTransaction(func(txRepo repo.User) (err error) {
		err = txRepo.CreateUser(ctx, 1, "Vasya")
		if err != nil {
			log.Fatalf("could not get user: %s\n", name)
		}

		// get user created in TX
		name, err = txRepo.GetUserNameByID(ctx, 1)
		fmt.Printf("user inside transaction: name '%s' err '%+v'\n", name, err)

		// no access before commit
		name, err = userRepo.GetUserNameByID(ctx, 1)
		fmt.Printf("user before commit: name '%s' err '%+v'\n", name, err)

		return nil
	})
	if err != nil {
		log.Fatalf("transaction: %s\n", err)
	}

	// user can be accessed from repo now
	name, err = userRepo.GetUserNameByID(ctx, 1)
	fmt.Printf("user after commit: name '%s' err '%+v'\n", name, err)

}

func rollbackTX(ctx context.Context, userRepo repo.User) {
	var name string
	var err = userRepo.DoTransaction(func(txRepo repo.User) (err error) {
		err = txRepo.CreateUser(ctx, 2, "Petya")
		if err != nil {
			log.Fatalf("could not get user: %s\n", name)
		}

		// get user created in TX
		name, err = txRepo.GetUserNameByID(ctx, 2)
		fmt.Printf("user inside transaction: name '%s' err '%+v'\n", name, err)

		// no access before commit
		name, err = userRepo.GetUserNameByID(ctx, 2)
		fmt.Printf("user before commit: name '%s' err '%+v'\n", name, err)

		return errors.New("here is rollback tx example")
	})
	if err != nil {
		fmt.Printf("tx rolled back: err '%+v'\n", err)
	}

	// no user here, cause of rollback
	name, err = userRepo.GetUserNameByID(ctx, 2)
	fmt.Printf("user after commit: name '%s' err '%+v'\n", name, err)

}
