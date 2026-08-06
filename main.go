package main

import (
	"context"
	"log"

	"transactions-lab/topics/transactions/challenges"
	database "transactions-lab/topics/transactions/database"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	// if err := challenges.One(ctx, db); err != nil {
	// 	log.Fatal(err)
	// }

	if err := challenges.Two(ctx, db, 1, 2, 100); err != nil {
		log.Fatal(err)
	}
	// sy.Run()
}
