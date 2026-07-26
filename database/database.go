package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

const databaseURL = "postgres://postgres:postgres@localhost:5432/transactions_lab?sslmode=disable"

func Connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL.")

	return db, nil
}
