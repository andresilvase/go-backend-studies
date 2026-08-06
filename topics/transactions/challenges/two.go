package challenges

import (
	"context"
	"database/sql"
	"fmt"

	customerrors "transactions-lab/topics/transactions/errors"
)

func Two(ctx context.Context, db *sql.DB, sourceWalletID, targetWalletID, amount int) error {
	fmt.Printf("Transfering money from wallet %d to wallet %d...\n", sourceWalletID, targetWalletID)

	if err := transferMoney(ctx, db, sourceWalletID, targetWalletID, amount); err != nil {
		return err
	}

	fmt.Println("Transference completed successfully!")
	return nil
}

func transferMoney(ctx context.Context, db *sql.DB, sourceWalletID, targetWalletID, amount int) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error starting transaction for transferMoney: %v", err),
		}
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND (balance - $1) >= 0`, amount, sourceWalletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error withdrawing money from account %d: %v", sourceWalletID, err),
		}
	}

	affectRows, err := result.RowsAffected()

	if err != nil || affectRows == 0 {
		return fmt.Errorf("there is no sufficient balance available or account %d doesn't exist\n", sourceWalletID)
	}

	result, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, targetWalletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error depositing account %d: %v", targetWalletID, err),
		}
	}

	affectRows, err = result.RowsAffected()

	if err != nil || affectRows == 0 {
		return fmt.Errorf("error ocurred when depositing money into wallet %d: %w\n", targetWalletID, err)
	}

	return tx.Commit()
}
