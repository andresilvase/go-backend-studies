package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const databaseURL = "postgres://postgres:postgres@localhost:5432/transactions_lab?sslmode=disable"

func Connect(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	fmt.Println("Connected to PostgreSQL.")

	return db, nil
}
