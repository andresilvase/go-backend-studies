package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"transactions-lab/topics/transactions/challenges"
	database "transactions-lab/topics/transactions/database"
	customeerors "transactions-lab/topics/transactions/errors"

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

	var dbErrors *customeerors.DBErr
	if err := challenges.Two(ctx, db, 1, 1, 100); err != nil {
		if errors.As(err, &dbErrors) {
			log.Fatalf("fatal error accessing DB: %v", dbErrors.Message)
		} else {
			fmt.Printf("%v", err)
		}
	}

	// sy.Run()
}
