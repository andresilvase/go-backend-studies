package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	db_conn "transactions-lab/topics/transactions/database"
	"transactions-lab/topics/transactions/database/pgstore"
	"transactions-lab/topics/utils"
)

func main() {
	ctx := context.Background()

	db, err := db_conn.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()
	queries := pgstore.New(db)

	var wg sync.WaitGroup

	bootstrap := []func(ctx context.Context, queries *pgstore.Queries){
		insertUsers,
		insertProducts,
	}

	wg.Add(len(bootstrap))

	for _, operation := range bootstrap {
		go func() {
			defer wg.Done()
			operation(ctx, queries)
		}()
	}

	wg.Wait()
}

func insertUsers(ctx context.Context, queries *pgstore.Queries) {

	fmt.Println("Inserindo usuários...")

	csvFolderPath, err := utils.GetCSVFolderPath()

	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Open(fmt.Sprintf("%s/%s.csv", csvFolderPath, utils.Users))

	if err != nil {
		log.Fatal(err)
	}

	defer csvFile.Close()

	fileReader := csv.NewReader(csvFile)

	csvRows, err := fileReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	for _, record := range csvRows[1:] {
		err := queries.CreateUser(ctx, record[1])

		if err != nil {
			log.Printf("failed to create user: %v", err)
		}
	}
}

func insertProducts(ctx context.Context, queries *pgstore.Queries) {
	fmt.Println("Inserindo produtos...")

	csvFolderPath, err := utils.GetCSVFolderPath()

	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Open(fmt.Sprintf("%s/%s.csv", csvFolderPath, utils.Products))

	if err != nil {
		log.Fatal(err)
	}

	defer csvFile.Close()

	fileReader := csv.NewReader(csvFile)

	csvRows, err := fileReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	for _, record := range csvRows[1:] {
		productPrice, err := strconv.ParseInt(record[2], 10, 64)

		if err != nil {
			log.Printf("invalid product price %q: %v", record[2], err)
			continue
		}

		newProductId, err := queries.CreateProduct(ctx, pgstore.CreateProductParams{
			Name:  record[1],
			Price: productPrice,
		})

		if err != nil {
			log.Printf("failed to create product: %v", err)
			continue
		}

		err = setProductInventory(ctx, queries, newProductId)

		if err != nil {
			log.Printf("failed to update product inventory: %v", err)
			continue
		}
	}
}

func setProductInventory(ctx context.Context, queries *pgstore.Queries, productId int64) error {
	return queries.SetProductInventory(
		ctx, pgstore.SetProductInventoryParams{
			ProductID: productId,
			// Stock:     9999,
			Stock: rand.Int63n(25) + 1,
		},
	)
}
