package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	sy "transactions-lab/topics/syntax"
	"transactions-lab/topics/transactions/challenges"
	mdl "transactions-lab/topics/transactions/challenges/models"
	database "transactions-lab/topics/transactions/database"
	customerrors "transactions-lab/topics/transactions/errors"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	// Essential config for challenge nine to thrive if it has nearly to 100 routines running.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	// runSyntax()
	// runChallengeOne(ctx, db)
	// runChallengeTwo(ctx, db)
	// runChallengeThree(ctx, db)
	// runChallengeFour(ctx, db)
	// runChallengeFive(ctx, db)
	// runChallengeSix(ctx, db)
	// runChallengeSeven(ctx, db)
	// runChallengeEight(ctx, db)
	runChallengeNine(ctx, db)
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
	var chTwoDbErrors *customerrors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 100

	var chTwoParam mdl.TransferParams = mdl.TransferParams{
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
	var chThreeDbErrors *customerrors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var simulatedFail = true
	var amount int64 = 250

	var chThreeParam mdl.TransferParams = mdl.TransferParams{
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
	var chFourDbErrors *customerrors.DBErr
	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var simulatedFail = true
	var amount int64 = 250

	var chFourParam mdl.TransferParams = mdl.TransferParams{
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

	var chFiveDbErrors *customerrors.DBErr
	var chFiveErrResult *customerrors.ErrResult

	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 250

	var chFiveParam mdl.TransferParams = mdl.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}

	fmt.Println("Starting Read Committed process...")

	if err := challenges.Five(chFiveParam, &wg); err != nil {
		if errors.As(err, &chFiveDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %s - %w", chFiveDbErrors.Message, chFiveDbErrors.Err))
		} else if errors.As(err, &chFiveErrResult) {
			if chFiveErrResult != nil {
				fmt.Println(chFiveErrResult.Error())
			}
		} else {
			fmt.Printf("error type: %T\n", err)
		}
	}

	wg.Wait()

	fmt.Println("Read Committed process finished!")
}

func runChallengeSix(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup

	var chSixDbErrors *customerrors.DBErr
	var chSixErrResult *customerrors.ErrResult

	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 250

	var chSixParam mdl.TransferParams = mdl.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}

	fmt.Println("Starting Repeatable Read process...")

	if err := challenges.Six(chSixParam, &wg); err != nil {
		if errors.As(err, &chSixDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %s - %w", chSixDbErrors.Message, chSixDbErrors.Err))
		} else if errors.As(err, &chSixErrResult) {
			if chSixErrResult != nil {
				fmt.Println(chSixErrResult.Error())
			}
		} else {
			fmt.Printf("error type: %T\n", err)
		}
	}

	wg.Wait()

	fmt.Println("Repeatable Read process finished!")
}

func runChallengeSeven(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup

	var chSevenDbErrors *customerrors.DBErr
	var chSevenErrResult *customerrors.ErrResult

	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 250

	var chSevenParam mdl.TransferParams = mdl.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}

	fmt.Println("Starting Phantom Reads process...")

	if err := challenges.Seven(chSevenParam, &wg); err != nil {
		if errors.As(err, &chSevenDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %s - %w", chSevenDbErrors.Message, chSevenDbErrors.Err))
		} else if errors.As(err, &chSevenErrResult) {
			if chSevenErrResult != nil {
				fmt.Println(chSevenErrResult.Error())
			}
		} else {
			fmt.Printf("error type: %T\n", err)
		}
	}

	wg.Wait()

	fmt.Println("Phantom Reads process finished!")
}

func runChallengeEight(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup

	var chEightDbErrors *customerrors.DBErr
	var chEightErrResult *customerrors.ErrResult

	var sourceWalletID int64 = 1
	var targetWalletID int64 = 2
	var amount int64 = 250

	var chEightParam mdl.TransferParams = mdl.TransferParams{
		Ctx:            &ctx,
		DB:             db,
		SourceWalletID: &sourceWalletID,
		TargetWalletID: &targetWalletID,
		Amount:         &amount,
	}

	fmt.Println("Starting Serializable process...")

	if err := challenges.Eight(chEightParam, &wg); err != nil {
		if errors.As(err, &chEightDbErrors) {
			log.Fatal(fmt.Errorf("fatal error accessing DB: %s - %w", chEightDbErrors.Message, chEightDbErrors.Err))
		} else if errors.As(err, &chEightErrResult) {
			if chEightErrResult != nil {
				fmt.Println(chEightErrResult.Error())
			}
		} else {
			fmt.Printf("error type: %T\n", err)
		}
	}

	wg.Wait()

	fmt.Println("Serializable process finished!")
}

func runChallengeNine(ctx context.Context, db *sql.DB) {
	var wg sync.WaitGroup
	var chNineDbErrors *customerrors.DBErr
	var chNineErrResult *customerrors.ErrResult

	err := challenges.Nine(challenges.ChallengeNineParams{
		Ctx: ctx,
		DB:  db,
		WG:  &wg,
	})

	if err != nil {
		if errors.As(err, &chNineDbErrors) {
			log.Fatal(err)
		} else if errors.As(err, &chNineErrResult) {
			if chNineErrResult != nil {
				fmt.Println(chNineErrResult.Error())
			}

		} else {
			fmt.Printf("error type: %T\n", err)
		}
	}

	fmt.Println("Multiple table challenge finished.")

	wg.Wait()
}
