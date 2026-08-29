package errors

import "fmt"

type DBErr struct {
	Message string
	Err     error
}

func (d *DBErr) Error() string {
	return fmt.Sprintf("an error occured accessing database: %s - %v", d.Message, d.Err)
}

func (d *DBErr) Unwrap() error {
	return d.Err
}

type OperationErr struct {
	Message string
	Err     error
}

func (d *OperationErr) Error() string {
	return fmt.Sprintf("operational error: %s", d.Message)
}

type ErrResult struct {
	TxName string
	Err    error
}

func (e *ErrResult) Error() string {
	return fmt.Sprintf("an error ocurred on transaction %s: %v", e.TxName, e.Err)
}
