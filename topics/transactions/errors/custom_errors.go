package errors

import "fmt"

type DBErr struct {
	Message string
	Err     error
}

func (d *DBErr) Error() string {
	return fmt.Sprintf("an error occured accessing database %s %v", d.Message, d.Err)
}

func (d *DBErr) Unwrap() error {
	return d.Err
}

type OperationErr struct {
	Message string
}

func (d *OperationErr) Error() string {
	return fmt.Sprintf("operational error: %s", d.Message)
}
