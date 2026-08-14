package challenges

import (
	"fmt"

	structs "transactions-lab/topics/transactions/challenges/structs"
	customerrors "transactions-lab/topics/transactions/errors"
)

func Two(param structs.TransferParams) error {
	var (
		sourceWalletID = *param.SourceWalletID
		targetWalletID = *param.TargetWalletID
		amount         = *param.Amount
	)

	fmt.Printf("Transferring money from wallet %d to wallet %d...\n", sourceWalletID, targetWalletID)

	if amount <= 0 {
		return &customerrors.OperationErr{Message: "amount should be greater than zero"}
	}

	if sourceWalletID == targetWalletID {
		return &customerrors.OperationErr{Message: "source and target wallets must differ"}
	}

	if err := transferMoney(param); err != nil {
		return err
	}

	fmt.Println("Transfer completed successfully!")
	return nil
}

func transferMoney(param structs.TransferParams) error {
	var (
		db             = param.DB
		ctx            = *param.Ctx
		sourceWalletID = *param.SourceWalletID
		targetWalletID = *param.TargetWalletID
		amount         = *param.Amount
	)

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
			Message: fmt.Sprintf("error withdrawing money from account %d:", sourceWalletID),
			Err:     err,
		}
	}

	affectRows, err := result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error reading affected rows when withdrawing account %d:", sourceWalletID),
			Err:     err,
		}
	}

	if affectRows == 0 {
		return fmt.Errorf("there is no sufficient balance available or account %d doesn't exist", sourceWalletID)
	}

	result, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, targetWalletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error depositing account %d:", targetWalletID),
			Err:     err,
		}
	}

	affectRows, err = result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error occurred when depositing money into wallet %d:", targetWalletID),
			Err:     err,
		}
	}

	if affectRows == 0 {
		return fmt.Errorf("target wallet %d does not exist", targetWalletID)
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing money transfer",
			Err:     err,
		}
	}

	return nil
}
