package challenges

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	mdl "transactions-lab/topics/transactions/challenges/models"
	customerrors "transactions-lab/topics/transactions/errors"
)

func Six(param mdl.TransferParams, wg *sync.WaitGroup) error {
	errorChan := make(chan *customerrors.ErrResult, 2)
	txChan := make(chan struct{})

	wg.Add(2)

	go func() {
		var txName = "A"
		defer wg.Done()
		if err := transactionA_RR(param, txName, txChan); err != nil {
			errorChan <- &customerrors.ErrResult{
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
		if err := transactionB_RR(param, txName, txChan); err != nil {
			errorChan <- &customerrors.ErrResult{
				TxName: txName,
				Err:    err,
			}
		} else {
			errorChan <- nil
		}
	}()

	var firstErr error

	for i := 0; i < 2; i++ {
		err := <-errorChan

		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	close(errorChan)

	return firstErr
}

func transactionA_RR(param mdl.TransferParams, txName string, ch chan struct{}) error {

	var (
		db             = param.DB
		ctx            = *param.Ctx
		targetWalletID = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})

	if err != nil {
		ch <- struct{}{}
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Six - transactionA_RR",
			Err:     err,
		}
	}

	defer tx.Rollback()

	if err = readAndPrintAmount_RR(ctx, tx, targetWalletID, txName); err != nil {
		ch <- struct{}{}
		return err
	}

	ch <- struct{}{}

	<-ch

	if err = readAndPrintAmount_RR(ctx, tx, targetWalletID, txName); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Six - transactionA_RR",
			Err:     err,
		}
	}

	return nil
}

func transactionB_RR(param mdl.TransferParams, txName string, ch chan struct{}) error {
	<-ch

	var (
		db             = param.DB
		ctx            = *param.Ctx
		targetWalletID = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		ch <- struct{}{}
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Six - transactionB_RR",
			Err:     err,
		}
	}

	defer tx.Rollback()

	if err = updateAmount_RR(ctx, tx, targetWalletID); err != nil {
		ch <- struct{}{}
		return err
	}

	if err := tx.Commit(); err != nil {
		ch <- struct{}{}
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Six - transactionB_RR",
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s updated amount\n", txName)

	ch <- struct{}{}

	return nil
}

func readAndPrintAmount_RR(ctx context.Context, tx *sql.Tx, walletID int64, txName string) error {
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

func updateAmount_RR(ctx context.Context, tx *sql.Tx, walletID int64) error {

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance - 1 WHERE id = $1`, walletID,
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
