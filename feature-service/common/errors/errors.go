package errors

const (
	SUCCESS        = 200
	BAD_REQUEST    = 400
	UNAUTHORIZED   = 401
	FORBIDDEN      = 403
	NOT_FOUND      = 404
	INTERNAL_ERROR = 500
)

type BaseError interface {
	Error() string
	Code() int
}

type baseError struct {
	code int
	err  error
}

func (e *baseError) Error() string {
	return e.err.Error()
}

func (e *baseError) Code() int {
	return e.code
}

func NewBaseError(code int, err error) BaseError {
	return &baseError{
		code: code,
		err:  err,
	}
}
