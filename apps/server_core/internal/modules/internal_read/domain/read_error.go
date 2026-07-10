package domain

import "errors"

type ReadErrorCode string

const (
	ReadErrorSourceUnavailable ReadErrorCode = "source_unavailable"
	ReadErrorUnsupportedQuery  ReadErrorCode = "unsupported_query"
)

type ReadError struct {
	Code    ReadErrorCode
	Message string
	Cause   error
}

func (e *ReadError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}

func (e *ReadError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewReadError(code ReadErrorCode, message string, cause error) error {
	return &ReadError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func IsReadErrorCode(err error, code ReadErrorCode) bool {
	var readErr *ReadError
	if !errors.As(err, &readErr) {
		return false
	}
	return readErr.Code == code
}
