package challenges

import (
	"log"
	"math/rand"
	"strconv"
	shared "transactions-lab/topics/transactions/challenges/shared"
	"transactions-lab/topics/transactions/database/pgstore"
	"transactions-lab/topics/utils"
)

func GenerateOrderIntent() (OrderIntent, error) {
	productsCount := (rand.Intn(8) + 1)
	randomUser := rand.Intn(5)
	productMap := make(map[pgstore.Product]int64, productsCount)

	users, err := shared.GetUsersFromCSVData()

	if err != nil {
		return OrderIntent{}, nil
	}

	products, err := getProductsFromData()

	if err != nil {
		return OrderIntent{}, nil
	}

	for i := 0; i < productsCount; i++ {
		productMap[products[rand.Intn(len(products))]] = (rand.Int63n(25) + 1)
	}

	return OrderIntent{
		Buyer:    users[randomUser],
		Products: productMap,
	}, nil
}

func getProductsFromData() ([]pgstore.Product, error) {
	content, err := utils.CsvContent(string(utils.Products))

	if err != nil {
		return []pgstore.Product{}, nil
	}

	var products []pgstore.Product

	for _, record := range content[1:] {
		productID, err := strconv.ParseInt(record[0], 10, 64)

		if err != nil {
			log.Printf("invalid product id %q: %v", record[2], err)
			continue
		}

		productPrice, err := strconv.ParseInt(record[2], 10, 64)

		if err != nil {
			log.Printf("invalid product price %q: %v", record[2], err)
			continue
		}

		products = append(products, pgstore.Product{
			ID:    productID,
			Name:  record[1],
			Price: productPrice,
		})
	}

	return products, nil
}
