package challenges

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	mdl "transactions-lab/topics/transactions/challenges/models"
	customerrors "transactions-lab/topics/transactions/errors"
)

func Five(param mdl.TransferParams, wg *sync.WaitGroup) error {
	errorChan := make(chan *mdl.ErrResult, 2)
	txChan := make(chan struct{})

	wg.Add(2)

	go func() {
		var txName = "A"
		defer wg.Done()
		if err := transactionA(param, txName, txChan); err != nil {
			errorChan <- &mdl.ErrResult{
				TxName: txName,
				Err:    err,
			}
		} else {
			errorChan <- nil
		}
	}()

	go func() {
		var txName = "B"

		defer wg.Done()
		if err := transactionB(param, txName, txChan); err != nil {
			errorChan <- &mdl.ErrResult{
				TxName: txName,
				Err:    err,
			}
		} else {
			errorChan <- nil
		}
	}()

	return <-errorChan
}

func transactionA(param mdl.TransferParams, txName string, ch chan struct{}) error {

	var (
		db             = param.DB
		ctx            = *param.Ctx
		targetWalletID = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Five - transactionA",
			Err:     err,
		}
	}

	defer tx.Rollback()

	if err = readAndPrintAmount(ctx, tx, targetWalletID, txName); err != nil {
		return err
	}

	ch <- struct{}{}

	<-ch

	if err = readAndPrintAmount(ctx, tx, targetWalletID, txName); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Five - transactionA",
			Err:     err,
		}
	}

	return nil
}

func transactionB(param mdl.TransferParams, txName string, ch chan struct{}) error {
	<-ch

	var (
		db             = param.DB
		ctx            = *param.Ctx
		targetWalletID = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Five - transactionB",
			Err:     err,
		}
	}

	defer tx.Rollback()

	if err = updateAmount(ctx, tx, targetWalletID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Five - transactionB",
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s updated amount\n", txName)

	ch <- struct{}{}

	return nil
}

func readAndPrintAmount(ctx context.Context, tx *sql.Tx, walletID int64, txName string) error {
	var balance int64
	err := tx.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1`, walletID,
	).Scan(&balance)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error reading amount from account %d:", walletID),
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s read amount = %d\n", txName, balance)

	return nil
}

func updateAmount(ctx context.Context, tx *sql.Tx, walletID int64) error {

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = 2255 WHERE id = $1`, walletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error updating balance for walltet %d:", walletID),
			Err:     err,
		}
	}

	affectRows, err := result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error reading rows affected after updating balance for wallet %d ", walletID),
			Err:     err,
		}
	}

	if affectRows == 0 {
		return fmt.Errorf("wallet %d does not exist", walletID)
	}

	return nil
}
