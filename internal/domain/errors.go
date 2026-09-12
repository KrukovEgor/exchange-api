package domain

import (
	"errors"
	"fmt"
)

type ErrCode string

const (
	CodeInvalidInput        ErrCode = "INVALID_INPUT"
	CodeNotFound            ErrCode = "NOT_FOUND"
	CodeCacheMiss           ErrCode = "CACHE_MISS"
	CodeCacheUnavailable    ErrCode = "CACHE_UNAVAILABLE"
	CodeProviderInvalidData ErrCode = "PROVIDER_INVALID_DATA"
	CodeProviderUnavailable ErrCode = "PROVIDER_UNAVAILABLE"
	CodeProviderRejected    ErrCode = "PROVIDER_REJECTED"
	CodeInternal            ErrCode = "INTERNAL"
)

type Error struct {
	Code ErrCode
	Msg  string
	Op   string
	Err  error
}

func newError(code ErrCode, msg, op string, err error) *Error {
	return &Error{Code: code, Msg: msg, Op: op, Err: err}
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Op != "" && e.Err != nil:
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Msg, e.Err)
	case e.Op != "":
		return fmt.Sprintf("%s: %s", e.Op, e.Msg)
	case e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	default:
		return e.Msg
	}
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return false
	}

	var t *Error

	if !errors.As(target, &t) || t == nil {
		return false
	}

	return e.Code == t.Code
}

var (
	ErrInvalidInput        = &Error{Code: CodeInvalidInput}
	ErrNotFound            = &Error{Code: CodeNotFound}
	ErrCacheMiss           = &Error{Code: CodeCacheMiss}
	ErrCacheUnavailable    = &Error{Code: CodeCacheUnavailable}
	ErrProviderInvalidData = &Error{Code: CodeProviderInvalidData}
	ErrProviderUnavailable = &Error{Code: CodeProviderUnavailable}
	ErrProviderRejected    = &Error{Code: CodeProviderRejected}
	ErrInternal            = &Error{Code: CodeInternal}
)

func NewInvalidInputError(op, msg string, err error) *Error {
	return newError(CodeInvalidInput, op, msg, err)
}

func NewNotFoundError(op, msg string, err error) *Error {
	return newError(CodeNotFound, op, msg, err)
}

func NewCacheMissError(op, msg string, err error) *Error {
	return newError(CodeCacheMiss, op, msg, err)
}

func NewCacheUnavailableError(op, msg string, err error) *Error {
	return newError(CodeCacheUnavailable, op, msg, err)
}

func NewProviderInvalidDataError(op, msg string, err error) *Error {
	return newError(CodeProviderInvalidData, op, msg, err)
}

func NewProviderUnavailableError(op, msg string, err error) *Error {
	return newError(CodeProviderUnavailable, op, msg, err)
}

func NewProviderRejectedError(op, msg string, err error) *Error {
	return newError(CodeProviderRejected, op, msg, err)
}

func NewInternalError(op, msg string, err error) *Error {
	return newError(CodeInternal, op, msg, err)
}
