package main

import (
	"context"
	"log"

	database "transactions-lab/database"
	sy "transactions-lab/syntax"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	// challenges.One(ctx, db)
	// challenges.Two(ctx, db, 1, 2, 100)
	sy.Run()
}
