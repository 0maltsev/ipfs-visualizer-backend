package handlers

import (
	"fmt"
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ClusterError struct {
	Op  string
	Msg string
	Err error
}

func (e *ClusterError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("adapter.%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("adapter.%s: %s", e.Op, e.Msg)
}

func (e *ClusterError) Unwrap() error { return e.Err }

func NewClusterError(op, msg string, err error) error {
	return &ClusterError{Op: op, Msg: msg, Err: err}
}

type NodeError struct {
	Op  string
	Msg string
	Err error
}

func (e *NodeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("adapter.%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("adapter.%s: %s", e.Op, e.Msg)
}

func (e *NodeError) Unwrap() error { return e.Err }

func NewNodeError(op, msg string, err error) error {
	return &NodeError{Op: op, Msg: msg, Err: err}
}

type RequestError struct {
	Op  string
	Msg string
	Err error
}

func (e *RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("request.%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("request.%s: %s", e.Op, e.Msg)
}

func (e *RequestError) Unwrap() error { return e.Err }

func NewRequestError(op, msg string, err error) error {
	return &RequestError{Op: op, Msg: msg, Err: err}
}

type ResponseError struct {
	Op  string
	Msg string
	Err error
}

func (e *ResponseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("response.%s: %s: %v", e.Op, e.Msg, e.Err)
	}
	return fmt.Sprintf("response.%s: %s", e.Op, e.Msg)
}

func (e *ResponseError) Unwrap() error { return e.Err }

func NewResponseError(op, msg string, err error) error {
	return &ResponseError{Op: op, Msg: msg, Err: err}
}
