package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/andrdru/sqltx/examples/sql2repos/repo"
)

const dsn = "host=localhost port=5432 user=pguser password=pgpassword dbname=pgdb sslmode=disable"

func main() {
	var ctx = context.Background()

	var db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("could not open db: %s\n", err)
	}

	defer func() {
		_ = db.Close()
	}()

	createTable(ctx, db)

	var userRepo = repo.NewUser(db)
	var addressRepo = repo.NewAddress(db)

	commitTX(ctx, userRepo, addressRepo)
}

func createTable(ctx context.Context, db *sql.DB) {
	const query = `CREATE TABLE users(
		"id" BIGSERIAL,		
		"name" TEXT);`

	var _, err = db.ExecContext(ctx, query)
	if err != nil {
		log.Fatalf("could not create table: %s\n", err)
	}

	const query2 = `CREATE TABLE addresses(
		"id" BIGSERIAL,		
		"street" TEXT);`

	_, err = db.ExecContext(ctx, query2)
	if err != nil {
		log.Fatalf("could not create table: %s\n", err)
	}
}

func commitTX(ctx context.Context, userRepo repo.User, addressRepo repo.Address) {
	var name string

	var err = userRepo.TX(func(txUser repo.User) (err error) {
		return addressRepo.TX(func(txAddress repo.Address) error {

			err = txUser.CreateUser(ctx, 1, "Vasya")
			if err != nil {
				log.Fatalf("could not create user: %s\n", err)
			}

			err = txAddress.CreateAddress(ctx, 1, "Main")
			if err != nil {
				log.Fatalf("could not create address: %s\n", err)
			}

			// get user created in TX
			name, err = txUser.GetUserNameByID(ctx, 1)
			fmt.Printf("user inside transaction: name '%s' err '%+v'\n", name, err)

			// get address created in TX
			name, err = txAddress.GetAddressStreetByID(ctx, 1)
			fmt.Printf("address inside transaction: name '%s' err '%+v'\n", name, err)

			// no access before commit
			name, err = userRepo.GetUserNameByID(ctx, 1)
			fmt.Printf("user before commit: name '%s' err '%+v'\n", name, err)

			// no access before commit
			name, err = addressRepo.GetAddressStreetByID(ctx, 1)
			fmt.Printf("address before commit: street '%s' err '%+v'\n", name, err)

			return nil
		})
	})

	if err != nil {
		log.Fatalf("transaction rollback: %s\n", err)
	}

	// user can be accessed from repo now
	name, err = userRepo.GetUserNameByID(ctx, 1)
	fmt.Printf("user after commit: name '%s' err '%+v'\n", name, err)

	// address can be accessed from repo now
	name, err = addressRepo.GetAddressStreetByID(ctx, 1)
	fmt.Printf("address after commit: street '%s' err '%+v'\n", name, err)

}
