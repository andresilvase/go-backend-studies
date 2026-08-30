package challenges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
	"transactions-lab/topics/transactions/database/pgstore"
	customerrors "transactions-lab/topics/transactions/errors"
	customerrs "transactions-lab/topics/transactions/errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type OrderIntent struct {
	Buyer    pgstore.User
	Products map[pgstore.Product]int64
}

type ChallengeNineParams struct {
	Ctx context.Context
	DB  *sql.DB
	WG  *sync.WaitGroup
}

func spawnGORoutines(orderIntents []OrderIntent, params ChallengeNineParams, errorChan chan *customerrs.ErrResult) {
	wg := params.WG
	wg.Add(len(orderIntents))

	for txNumber := 0; txNumber < len(orderIntents); txNumber++ {
		go func() {
			defer wg.Done()
			err := retryAttempts(func(attempt int) error {
				fmt.Printf("Running attempt %d for Tx-%d\n", attempt, txNumber)
				err := makePurchase(params, orderIntents[txNumber])

				return err
			}, txNumber)

			if err != nil {
				errorChan <- &customerrors.ErrResult{
					TxName: fmt.Sprintf("Tx-%d", txNumber),
					Err:    err,
				}
			} else {
				errorChan <- nil
			}
		}()
	}
}

func retryAttempts(operation func(attempt int) error, txNumber int) error {
	MAX_ATTEMPTS := 20
	var lastErr error
	initialBackoffDelay := 100 * time.Millisecond

	for attempt := 0; attempt < MAX_ATTEMPTS; attempt++ {
		fmt.Printf("\nTransaction %d, attempt %d\n", txNumber, attempt)
		lastErr = operation(attempt)

		if lastErr == nil {
			return nil
		}

		if !shouldRetry(lastErr) {
			fmt.Printf("Should not retry %v - %d\n", lastErr, attempt)
			return lastErr
		}

		if MAX_ATTEMPTS-1 == attempt {
			fmt.Printf("Should not retry anymore. Attempt %d, transaction %d\n", attempt, txNumber)
			return lastErr
		}

		delay := initialBackoffDelay * time.Duration(1<<attempt)
		time.Sleep(delay)
	}

	return nil
}

func shouldRetry(err error) bool {
	var serializableErr *pgconn.PgError

	if errors.As(err, &serializableErr) {
		if serializableErr.Code == "40001" || serializableErr.Code == "40P01" {
			return true
		}
	}

	return false
}

func getOrderIntents() ([]OrderIntent, error) {
	NUMBER_OF_ORDER_INTENTS := 500
	orderIntents := make([]OrderIntent, NUMBER_OF_ORDER_INTENTS)
	var firstErr error

	for i := 0; i < NUMBER_OF_ORDER_INTENTS; i++ {
		orderIntent, err := GenerateOrderIntent()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		} else {
			orderIntents[i] = orderIntent
		}
	}

	if firstErr != nil {
		return []OrderIntent{}, firstErr
	}

	return orderIntents, nil
}

func ninePrintTest() error {
	result, err := GenerateOrderIntent()

	if err != nil {
		return err
	}

	fmt.Printf("%v", result)

	return nil
}

func Nine(params ChallengeNineParams) error {
	fmt.Println("Running Challenge 9...")
	// ================ START OF TEST =================
	// _ := ninePrintTest()
	// ================ END OF TEST =================

	// ================ START OF PRODUCTION CODE =================

	orderIntents, err := getOrderIntents()

	if err != nil {
		return &customerrs.OperationErr{
			Message: "challenge Nine Error",
			Err:     err,
		}
	}

	errorChan := make(chan *customerrs.ErrResult, len(orderIntents))

	spawnGORoutines(orderIntents, params, errorChan)

	var firstErr error

	for i := 0; i < len(orderIntents); i++ {
		err := <-errorChan

		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
	// ================ END OF PRODUCTION CODE =================
}

func makePurchase(params ChallengeNineParams, orderIntent OrderIntent) error {
	var (
		ctx = params.Ctx
		db  = params.DB
	)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return &customerrs.DBErr{
			Message: "não foi possível iniciar transação",
			Err:     err,
		}
	}

	txQuery := pgstore.New(db).WithTx(tx)

	// Create an order
	orderId, err := txQuery.CreateOrder(ctx, orderIntent.Buyer.ID)

	for product, qtdRequested := range orderIntent.Products {
		// Read product stock
		stock, err := readProductStock(ctx, txQuery, product.ID)

		fmt.Printf("\nStock for %s: %d.\n", product.Name, stock)
		if err != nil {
			return &customerrs.DBErr{
				Message: fmt.Sprintf("error reading stock for product %s.", product.Name),
				Err:     err,
			}
		}

		// Validate quantity
		if stock < qtdRequested {
			return &customerrs.OperationErr{
				Message: fmt.Sprintf("the remaining stock of %s is shorter than the amount requested: there is %d. Requested %d.", product.Name, stock, qtdRequested),
			}
		}

		newInventory := stock - qtdRequested

		// Decrease inventory
		err = updateProductInventory(ctx, txQuery, product.ID, newInventory)
		if err != nil {
			return &customerrs.DBErr{
				Message: fmt.Sprintf("error decreasing inventory for product %s.", product.Name),
				Err:     err,
			}
		}

		fmt.Printf("New stock amount for %s: %d.\n", product.Name, newInventory)

		fmt.Println("Creating order item...")
		// Create order item
		err = txQuery.CreateOrderItem(ctx, pgstore.CreateOrderItemParams{
			OrderID:   orderId,
			ProductID: product.ID,
			Quantity:  qtdRequested,
		})

		if err != nil {
			return &customerrs.DBErr{
				Message: fmt.Sprintf("error inserting product %s into order items.", product.Name),
				Err:     err,
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return &customerrs.DBErr{
			Message: "error committing purchase transaction",
			Err:     err,
		}
	}

	return nil
}

func readProductStock(ctx context.Context, txQuery *pgstore.Queries, productID int64) (int64, error) {
	stock, err := txQuery.GetStockForProductID(ctx, productID)

	if err != nil {
		return 0, err
	}

	return stock, nil
}

func updateProductInventory(ctx context.Context, txQuery *pgstore.Queries, productID, newStock int64) error {
	return txQuery.UpdateProductInventory(ctx, pgstore.UpdateProductInventoryParams{
		ProductID: productID,
		Stock:     newStock,
	})
}
