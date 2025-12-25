package psql_connection

import "fmt"

type PostgresConnectionError struct {
	Op  string
	Msg string
	Err error
}

func (e *PostgresConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("psql_connection.%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("psql_connection.%s: %s", e.Op, e.Msg)
}

func (e *PostgresConnectionError) Unwrap() error { return e.Err }

func newPostgresConnectionError(op, msg string, err error) error {
	return &PostgresConnectionError{Op: op, Msg: msg, Err: err}
}
