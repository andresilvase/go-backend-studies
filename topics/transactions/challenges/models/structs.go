package models

import (
	"context"
	"database/sql"
)

type TransferParams struct {
	Ctx            *context.Context
	DB             *sql.DB
	SourceWalletID *int64
	TargetWalletID *int64
	SimulatedFail  *bool
	Amount         *int64
}
