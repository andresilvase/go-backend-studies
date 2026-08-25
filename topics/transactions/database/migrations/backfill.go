package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"transactions-lab/topics/transactions/database"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	for i := 0; i < 5; i++ {
		if _, err := db.Exec(
			"INSERT INTO products (name, price, quantity) VALUES($1, $2, $3)",
			fmt.Sprintf("Product-%d", i), (i+1)*(i+1<<1), math.Pow(float64(i+3), 2),
		); err != nil {
			log.Fatal(err)
		}
	}
}
