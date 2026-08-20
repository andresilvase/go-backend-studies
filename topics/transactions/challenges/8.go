package challenges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
	mdl "transactions-lab/topics/transactions/challenges/models"
	customerrors "transactions-lab/topics/transactions/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func Eight(param mdl.TransferParams, wg *sync.WaitGroup) error {
	ctx := *param.Ctx
	CHAN_BUF_SIZE := 2

	errorChan := make(chan *customerrors.ErrResult, CHAN_BUF_SIZE)

	// Sync channels must have at leat one size in buffer because without it,
	// when one send there will be no one to reveive and will fall into a deadlock error.
	aReady := make(chan struct{}, 1)
	bReady := make(chan struct{}, 1)

	// Context used to cancel the experiment in case one of the goroutines fails before the sync point.
	barrierCtx, cancelCtx := context.WithCancel(ctx)

	wg.Add(2)

	afterRead := func(signalingCh, waitingCh chan struct{}) error {
		signalingCh <- struct{}{}

		select {
		case <-waitingCh:
			return nil
		case <-barrierCtx.Done():
			return barrierCtx.Err()
		}
	}

	go func() {
		var txName = "A"
		defer wg.Done()

		err := retryWrapper(txName, func(attempt int) error {
			reachedBarrier := false
			if attempt == 0 {
				err := transactionASerial(barrierCtx, param, txName, func() error {
					reachedBarrier = true
					return afterRead(aReady, bReady)
				})

				if err != nil && !reachedBarrier {
					cancelCtx()
				}

				return err
			} else {
				return transactionASerial(ctx, param, txName, nil)
			}
		})

		if err != nil {
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

		err := retryWrapper(txName, func(attempt int) error {
			reachedBarrier := false
			if attempt == 0 {
				err := transactionBSerial(barrierCtx, param, txName, func() error {
					reachedBarrier = true
					return afterRead(bReady, aReady)
				})
				if err != nil && !reachedBarrier {
					cancelCtx()
				}
				return err
			} else {
				return transactionBSerial(ctx, param, txName, nil)
			}
		})

		if err != nil {
			errorChan <- &customerrors.ErrResult{
				TxName: txName,
				Err:    err,
			}
		} else {
			errorChan <- nil
		}
	}()

	var firstErr error

	for i := 0; i < CHAN_BUF_SIZE; i++ {
		err := <-errorChan

		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	close(errorChan)

	return firstErr
}

func retryWrapper(txName string, operation func(attempt int) error) error {
	MAX_ATTEMPT_RETRY := 5
	initialDelay := 100 * time.Millisecond

	for attempt := 0; attempt < MAX_ATTEMPT_RETRY; attempt++ {
		fmt.Printf("Tx-%s running attempt %d...\n", txName, attempt)
		err := operation(attempt)

		if err == nil {
			return nil
		}

		if !isRetryable(err) {
			return err
		}

		delay := initialDelay * time.Duration(1<<attempt)
		fmt.Printf("Tx-%s retrying after %s\n", txName, delay)
		time.Sleep(delay)
	}

	return nil
}

func isRetryable(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "40001" || pgErr.Code == "40P01" {
			return true
		}
	}

	return false
}

func transactionASerial(ctx context.Context, param mdl.TransferParams, txName string, afterReadA func() error) error {

	var (
		db           = param.DB
		sourceWallet = *param.SourceWalletID
		targetWallet = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Eight - transactionASerial",
			Err:     err,
		}
	}

	defer tx.Rollback()

	targetWalletBalance, err := readOtherWalletBalance(ctx, tx, txName, targetWallet)

	if err != nil {
		return err
	}

	if afterReadA != nil {
		if err := afterReadA(); err != nil {
			return err
		}
	}

	if targetWalletBalance >= 600 {
		if err = zeroThisWalletBalance(ctx, tx, txName, sourceWallet); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Eight - transactionASerial",
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s finished\n", txName)

	return nil
}

func transactionBSerial(ctx context.Context, param mdl.TransferParams, txName string, afterReadB func() error) error {

	var (
		db           = param.DB
		targetWallet = *param.SourceWalletID
		sourceWallet = *param.TargetWalletID
	)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return &customerrors.DBErr{
			Message: "error starting transaction for challenge Eight - transactionBSerial",
			Err:     err,
		}
	}

	defer tx.Rollback()

	targetWalletBalance, err := readOtherWalletBalance(ctx, tx, "B", targetWallet)

	if err != nil {
		return err
	}

	if afterReadB != nil {
		if err := afterReadB(); err != nil {
			return err
		}
	}

	if targetWalletBalance >= 600 {
		if err := zeroThisWalletBalance(ctx, tx, txName, sourceWallet); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return &customerrors.DBErr{
			Message: "error committing transaction for challenge Eight - transactionBSerial",
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s finished\n", txName)

	return nil
}

func readOtherWalletBalance(ctx context.Context, tx *sql.Tx, txName string, walletID int64) (int64, error) {
	var targetWalletBalance int64

	err := tx.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1;`,
		walletID,
	).Scan(&targetWalletBalance)

	if err != nil {
		return 0, &customerrors.DBErr{
			Message: fmt.Sprintf("error counting predicate in transferMoneyToAccountB for %s: %s", txName, err.Error()),
			Err:     err,
		}
	}

	fmt.Printf("Tx-%s read wallet %d balance = %d\n", txName, walletID, targetWalletBalance)

	return targetWalletBalance, nil
}

func zeroThisWalletBalance(ctx context.Context, tx *sql.Tx, txName string, walletID int64) error {
	functionName := "zeroThisWalletBalance"

	result, err := tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = 0 WHERE id = $1;`,
		walletID,
	)

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error zeroing wallet %d balance in %s for %s: %s", walletID, functionName, txName, err.Error()),
			Err:     err,
		}
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		return &customerrors.DBErr{
			Message: fmt.Sprintf("error reading affectedRows in %s for %s: %s", functionName, txName, err.Error()),
			Err:     err,
		}
	}

	if affectedRows == 0 {
		return &customerrors.OperationErr{
			Message: fmt.Sprintf("seems like a wallet under ID %d doesn't exist in %s for %s.", walletID, functionName, txName),
		}
	}

	return nil
}
