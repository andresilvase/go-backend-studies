package models

import (
	"context"
	"database/sql"
	"fmt"
)

type TransferParams struct {
	Ctx            *context.Context
	DB             *sql.DB
	SourceWalletID *int64
	TargetWalletID *int64
	SimulatedFail  *bool
	Amount         *int64
}

type ErrResult struct {
	TxName string
	Err    error
}

func (e *ErrResult) Error() string {
	return fmt.Sprintf("error ocurred on transaction %s: %v", e.TxName, e.Err)
}
