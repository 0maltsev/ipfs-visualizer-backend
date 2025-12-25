package app

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	KindClient     ErrorKind = "client"
	KindConfig     ErrorKind = "config"
	KindStorage    ErrorKind = "storage"
	KindIO         ErrorKind = "io"
	KindValidation ErrorKind = "validation"
	KindNotFound   ErrorKind = "not_found"
	KindServer     ErrorKind = "server"
	KindUnknown    ErrorKind = "unknown"
)

type AppError struct {
	Op   string
	Kind ErrorKind
	Msg  string
	Err  error
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Msg)
}

func (e *AppError) Unwrap() error { return e.Err }

func NewClientError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindClient, Msg: msg, Err: err}
}

func NewConfigError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindConfig, Msg: msg, Err: err}
}

func NewStorageError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindStorage, Msg: msg, Err: err}
}

func NewIOError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindIO, Msg: msg, Err: err}
}

func NewValidationError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindValidation, Msg: msg, Err: err}
}

func NewNotFoundError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindNotFound, Msg: msg, Err: err}
}

func NewServerError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindServer, Msg: msg, Err: err}
}

func NewUnknownError(op, msg string, err error) error {
	return &AppError{Op: op, Kind: KindUnknown, Msg: msg, Err: err}
}

func IsKind(err error, kind ErrorKind) bool {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Kind == kind
	}
	return false
}