package sqlmodelerrors

import "fmt"

type PostgresModelError struct {
	FuncName string
	Msg      string
	Err      error
}

func (e *PostgresModelError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.FuncName, e.Msg, e.Err)
}

func (e *PostgresModelError) Unwrap() error {
	return e.Err
}

func NewPostgresModelError(funcName, msg string, err error) *PostgresModelError {
	return &PostgresModelError{
		FuncName: funcName,
		Msg:      msg,
		Err:      err,
	}
}
