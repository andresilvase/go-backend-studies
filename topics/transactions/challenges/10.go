package challenges

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"transactions-lab/topics/transactions/database/pgstore"
	customerrors "transactions-lab/topics/transactions/errors"
)

type ChallengeTenParams struct {
	Ctx context.Context
	DB  *sql.DB
}

func Ten(params ChallengeTenParams) error {

	var (
		ctx = params.Ctx
		db  = params.DB
	)

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return &customerrors.DBErr{
			Message: "error initiating transaction for challenge ten.",
			Err:     err,
		}
	}

	defer tx.Rollback()

	txQuery := pgstore.New(db).WithTx(tx)

	userData, err := GetRandomUser()

	if err != nil {
		return &customerrors.OperationErr{
			Message: "error gettting random user",
			Err:     err,
		}
	}

	createdUserId, err := CreateUser(ctx, *txQuery, userData)

	if err != nil {
		return &customerrors.OperationErr{
			Message: "error creating user",
			Err:     err,
		}
	}

	createdWalletId, err := CreateWallet(ctx, *txQuery, createdUserId)

	if err != nil {
		return &customerrors.OperationErr{
			Message: fmt.Sprintf("error wallet for user user %s with ID %d", userData.Name, createdUserId),
			Err:     err,
		}
	}

	err = CreateProfile(ctx, *txQuery, createdWalletId, createdUserId)

	if err != nil {
		return &customerrors.OperationErr{
			Message: fmt.Sprintf("error profile for user user %s with ID %d", userData.Name, createdUserId),
			Err:     err,
		}
	}

	if err = tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error initiating transaction for challenge ten.",
			Err:     err,
		}
	}

	return nil
}

func CreateUser(ctx context.Context, txQuery pgstore.Queries, userData pgstore.User) (int64, error) {
	userId, err := txQuery.CreateUser(ctx, userData.Name)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func CreateWallet(ctx context.Context, txQuery pgstore.Queries, userId int64) (int64, error) {
	walletId, err := txQuery.CreateWallet(ctx, pgstore.CreateWalletParams{
		UserID:  userId,
		Balance: rand.Int63n(9999),
	})

	if err != nil {
		return 0, err
	}

	return walletId, nil
}

func CreateProfile(ctx context.Context, txQuery pgstore.Queries, walletId, userId int64) error {

	profileCategories := []string{
		"MASTER",
		"BREGA",
		"HIPER",
		"BLASTER",
	}

	if err := txQuery.CreateProfile(ctx, pgstore.CreateProfileParams{
		Category: profileCategories[rand.Intn(len(profileCategories))],
		WalletID: walletId,
		UserID:   userId,
	}); err != nil {
		return err
	}

	return nil
}
