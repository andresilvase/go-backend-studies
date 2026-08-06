package challenges

import (
	"context"
	"database/sql"
	"fmt"

	customerrors "transactions-lab/topics/transactions/errors"
)

func Two(ctx context.Context, db *sql.DB, sourceWalletID, targetWalletID, amount int) error {
	fmt.Printf("Transfering money from wallet %d to wallet %d...\n", sourceWalletID, targetWalletID)

	if amount <= 0 {
		return &customerrors.OperationErr{Message: "amount should be greater than zero"}
	}

	if sourceWalletID == targetWalletID {
		return &customerrors.OperationErr{Message: "source and target wallets must differ"}
	}

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
			Message: "error starting transaction for transferMoney",
			Err:     err,
		}
	}

	defer tx.Rollback()

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND (balance - $1) >= 0`, amount, sourceWalletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error withdrawing money from account %d:\n", sourceWalletID),
			Err:     err,
		}
	}

	affectRows, err := result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error reading affected rows when withdrawing account %d:\n", sourceWalletID),
			Err:     err,
		}
	}

	if affectRows == 0 {
		return fmt.Errorf("there is no sufficient balance available or account %d doesn't exist\n", sourceWalletID)
	}

	result, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, targetWalletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error depositing account %d:\n", targetWalletID),
			Err:     err,
		}
	}

	affectRows, err = result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error ocurred when depositing money into wallet %d:\n", targetWalletID),
			Err:     err,
		}
	}

	if affectRows == 0 {
		return fmt.Errorf("error ocurred when depositing money into wallet %d\n", targetWalletID)
	}

	return tx.Commit()
}
