package main

import (
	"context"
	"fmt"
	"log"

	challenges "transactions-lab/challenges"
	database "transactions-lab/database"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	challenges.One(ctx, db)

	// arr := []int{1, 2, 3, 6, 4, 8, 9, 6} // 4
	// arr := []int{1, 2, 1, 1, 1, 1, 1, 1} // 2
	// arr := []int{1, 1, 1, 1, 1, 1, 1, 1} // 1
	// arr := []int{9, 8, 7, 6, 5, 4, 3, 1} // 1
	arr := []int{9, 8, 7, 6, 5, 1, 2, 3} // 3

	output := challenges.StrictlyIncreasingArrayLength(arr)

	fmt.Printf("\n\nStrictly Increasing Array Length for %v is... %v.\n", arr, output)
}
