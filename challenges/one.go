package challenges

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func One(ctx context.Context, db *sql.DB) {

	userID, walletID, err := createUserWithWallet(ctx, User{Name: "John Doe"}, db)

	if err != nil {
		log.Fatalf("Error creating user and wallet:\n%v", err)
	}

	fmt.Printf("User created with ID: %+v\n", userID)
	fmt.Printf("Wallet created with ID: %+v\n", walletID)
}

func createUserWithWallet(ctx context.Context, user User, db *sql.DB) (int16, int16, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return -1, -1, err
	}

	defer tx.Rollback()

	userID, err := createUser(ctx, user, tx)

	if err != nil {
		return -1, -1, fmt.Errorf("Error creating user: %v", err)
	}

	log.Printf("User created with ID: %+v\n", userID)

	walletID, err := createWallet(ctx, Wallet{UserID: userID, Balance: 0}, tx)

	if err != nil {
		return -1, -1, err
	}

	log.Printf("Wallet created for user ID: %+v\n", userID)

	if err = tx.Commit(); err != nil {
		return -1, -1, fmt.Errorf("Error committing transaction: %w", err)
	}

	return userID, walletID, nil
}

func createUser(ctx context.Context, user User, tx *sql.Tx) (int16, error) {

	var userID int16

	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO users (name) VALUES ($1) returning id`,
		user.Name,
	).Scan(&userID)

	if err != nil {
		return -1, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func createWallet(ctx context.Context, wallet Wallet, tx *sql.Tx) (int16, error) {
	var walletID int16

	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO wallets (user_id, balance) VALUES ($1, $2) returning id`,
		wallet.UserID,
		wallet.Balance,
	).Scan(&walletID)

	if err != nil {
		return -1, fmt.Errorf("failed to create wallet: %w", err)
	}

	return walletID, nil
}

type Wallet struct {
	UserID  int16
	Balance int64
}

type User struct {
	ID   int16
	Name string
}
