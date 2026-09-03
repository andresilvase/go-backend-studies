package csv_data_shared

import (
	"log"
	"strconv"
	"transactions-lab/topics/transactions/database/pgstore"
	"transactions-lab/topics/utils"
)

func GetUsersFromCSVData() ([]pgstore.User, error) {
	content, err := utils.CsvContent(string(utils.Users))

	if err != nil {
		return []pgstore.User{}, nil
	}

	var users []pgstore.User

	for _, record := range content[1:] {
		userID, err := strconv.ParseInt(record[0], 10, 64)

		if err != nil {
			log.Printf("invalid user id %q: %v", record[2], err)
			continue
		}

		users = append(users, pgstore.User{
			ID:   userID,
			Name: record[1],
		})
	}

	return users, nil
}
