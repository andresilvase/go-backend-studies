package challenges

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func Two(ctx context.Context, db *sql.DB, sourceWalletID, targetWalletID, amount int) {
	fmt.Printf("Transfering money from wallet %d to wallet %d...\n", sourceWalletID, targetWalletID)

	if err := transferMoney(ctx, db, sourceWalletID, targetWalletID, amount); err != nil {
		log.Fatal("error ocurred when transfering money")
	}

	fmt.Println("Transference completed successfully!")
}

func transferMoney(ctx context.Context, db *sql.DB, sourceWalletID, targetWalletID, amount int) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("error when starting transaction for transference %w", err)
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance - $1 WHERE id = $2`, amount, sourceWalletID,
	)

	affectRows, _ := result.RowsAffected()

	if err != nil || affectRows == 0 {
		return fmt.Errorf("an error ocurred when withdrawing the money from account %d\n", sourceWalletID)
	}

	result, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, targetWalletID,
	)

	affectRows, _ = result.RowsAffected()

	if err != nil || affectRows == 0 {
		return fmt.Errorf("error ocurred when depositing money into wallet %d\n", targetWalletID)
	}

	return tx.Commit()
}
