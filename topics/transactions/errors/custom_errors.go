package errors

import "fmt"

type DBErr struct {
	Message string
}

func (d *DBErr) Error() string {
	return fmt.Sprintf("an error occured accessing database %s", d.Message)
}
