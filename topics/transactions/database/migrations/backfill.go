package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	db_conn "transactions-lab/topics/transactions/database"
	"transactions-lab/topics/transactions/database/pgstore"
	"transactions-lab/topics/utils"
)

type TableName string

const (
	Users    TableName = "users"
	Products TableName = "products"
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

	csvFolderPath, err := getCSVFolderPath()

	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Open(fmt.Sprintf("%s/%s.csv", csvFolderPath, Users))

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

	csvFolderPath, err := getCSVFolderPath()

	if err != nil {
		log.Fatal(err)
	}

	csvFile, err := os.Open(fmt.Sprintf("%s/%s.csv", csvFolderPath, Products))

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

		err = queries.CreateProduct(ctx, pgstore.CreateProductParams{
			Name:  record[1],
			Price: productPrice,
		})

		if err != nil {
			log.Printf("failed to create product: %v", err)
			continue
		}
	}
}

func getCSVFolderPath() (string, error) {
	root, err := utils.ProjectRoot()

	if err != nil {
		return "", err
	}

	csvFolderPath := filepath.Join(
		root,
		"topics",
		"transactions",
		"data",
	)

	return csvFolderPath, nil
}
