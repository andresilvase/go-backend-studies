package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	sy "transactions-lab/topics/syntax"
	"transactions-lab/topics/transactions/challenges"
	structs "transactions-lab/topics/transactions/challenges/structs"
	database "transactions-lab/topics/transactions/database"
	customeerors "transactions-lab/topics/transactions/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	// runSyntax()
	// runChallengeOne(ctx, db)
	// runChallengeTwo(ctx, db)
	// runChallengeThree(ctx, db)
	// runChallengeFour(ctx, db)
	runChallengeFive(ctx, db)
}

func runSyntax() {
	sy.Run()
}

func runChallengeOne(ctx context.Context, db *sql.DB) {
	if err := challenges.One(ctx, db); err != nil {
		log.Fatal(err)
	}
}

func runChallengeTwo(ctx context.Context, db *sql.DB) {
	var chTwoDbErrors *customeerors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 100

	var chTwoParam structs.TransferParams = structs.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}
	if err := challenges.Two(chTwoParam); err != nil {
		if errors.As(err, &chTwoDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %v - %w", chTwoDbErrors.Message, chTwoDbErrors.Err))
		} else {
			fmt.Println(err)
		}
	}
}

func runChallengeThree(ctx context.Context, db *sql.DB) {
	var chThreeDbErrors *customeerors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var simulatedFail = true
	var amount int64 = 250

	var chThreeParam structs.TransferParams = structs.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		SimulatedFail:  &simulatedFail,
		Amount:         &amount,
	}
	if err := challenges.Three(chThreeParam); err != nil {
		if errors.As(err, &chThreeDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %v - %w", chThreeDbErrors.Message, chThreeDbErrors.Err))
		} else {
			fmt.Println(err)
		}
	}
}

func runChallengeFour(ctx context.Context, db *sql.DB) {
	var chFourDbErrors *customeerors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var simulatedFail = true
	var amount int64 = 250

	var chFourParam structs.TransferParams = structs.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		SimulatedFail:  &simulatedFail,
		Amount:         &amount,
	}

	if err := challenges.Four(chFourParam); err != nil {
		if errors.As(err, &chFourDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %v - %w", chFourDbErrors.Message, chFourDbErrors.Err))
		} else {
			fmt.Println(err)
		}
	}
}

func runChallengeFive(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup

	var chFiveDbErrors *customeerors.DBErr

	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 250

	var chFiveParam structs.TransferParams = structs.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}

	fmt.Println("Starting Read Committed process...")

	if err := challenges.Five(chFiveParam, &wg); err != nil {
		if errors.As(err, &chFiveDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %v - %w", chFiveDbErrors.Message, chFiveDbErrors.Err))
		} else {
			fmt.Println(err)
		}
	}

	wg.Wait()

	fmt.Println("Read Committed process finished!")
}
