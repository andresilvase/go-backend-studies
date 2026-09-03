package challenges

import (
	"math/rand"
	shared "transactions-lab/topics/transactions/challenges/shared"
	"transactions-lab/topics/transactions/database/pgstore"
)

func GetRandomUser() (pgstore.User, error) {
	users, err := shared.GetUsersFromCSVData()

	if err != nil {
		return pgstore.User{}, err
	}

	return users[rand.Intn(len(users))], nil
}
