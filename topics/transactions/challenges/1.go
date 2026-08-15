package challenges

import (
	"context"
	"database/sql"
	"fmt"
)

func One(ctx context.Context, db *sql.DB) error {

	userID, walletID, err := createUserWithWallet(ctx, User{Name: "Shang Tsung"}, db)

	if err != nil {
		return err
	}

	fmt.Printf("User created with ID: %+v\n", userID)
	fmt.Printf("Wallet created with ID: %+v\n", walletID)
	return nil
}

func createUserWithWallet(ctx context.Context, user User, db *sql.DB) (int64, int64, error) {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return 0, 0, fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	userID, err := createUser(ctx, user, tx)

	if err != nil {
		return 0, 0, fmt.Errorf("error creating user: %w", err)
	}

	walletID, err := createWallet(ctx, Wallet{UserID: userID, Balance: 0}, tx)

	if err != nil {
		return 0, 0, fmt.Errorf("error creating wallet: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return userID, walletID, nil
}

func createUser(ctx context.Context, user User, tx *sql.Tx) (int64, error) {

	var userID int64

	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO users (name) VALUES ($1) returning id`,
		user.Name,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func createWallet(ctx context.Context, wallet Wallet, tx *sql.Tx) (int64, error) {
	var walletID int64

	err := tx.QueryRowContext(
		ctx,
		`INSERT INTO wallets (user_id, balance) VALUES ($1, $2) returning id`,
		wallet.UserID,
		wallet.Balance,
	).Scan(&walletID)

	if err != nil {
		return 0, fmt.Errorf("failed to create wallet: %w", err)
	}

	return walletID, nil
}

type Wallet struct {
	UserID  int64
	Balance int64
}

type User struct {
	ID   int64
	Name string
}
