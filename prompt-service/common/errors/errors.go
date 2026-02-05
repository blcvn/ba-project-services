package errors

import "fmt"

type ErrorCode int

const (
	BAD_REQUEST    ErrorCode = 400
	NOT_FOUND      ErrorCode = 404
	CONFLICT_ERROR ErrorCode = 409
	INTERNAL_ERROR ErrorCode = 500
)

type BaseError interface {
	Error() string
	GetCode() ErrorCode
}

type baseError struct {
	code ErrorCode
	err  error
}

func NewBaseError(code ErrorCode, err error) BaseError {
	return &baseError{code: code, err: err}
}

func (e *baseError) Error() string      { return e.err.Error() }
func (e *baseError) GetCode() ErrorCode { return e.code }

func BadRequest(msg string) BaseError { return NewBaseError(BAD_REQUEST, fmt.Errorf(msg)) }
func NotFound(msg string) BaseError   { return NewBaseError(NOT_FOUND, fmt.Errorf(msg)) }
func Conflict(msg string) BaseError   { return NewBaseError(CONFLICT_ERROR, fmt.Errorf(msg)) }
func Internal(err error) BaseError    { return NewBaseError(INTERNAL_ERROR, err) }
