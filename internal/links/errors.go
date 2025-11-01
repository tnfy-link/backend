package links

import (
	"errors"
)

var (
	ErrValidationFailed  = errors.New("validation failed")
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
)

type ValidationError struct {
	field string
	err   error
}

func (v ValidationError) Error() string {
	return v.field + ": " + v.err.Error()
}

func (v ValidationError) Unwrap() error {
	return v.err
}

func newValidationError(field string, err error) ValidationError {
	return ValidationError{
		field: field,
		err:   err,
	}
}

func AsValidationError(err error) *ValidationError {
	var v ValidationError

	if errors.As(err, &v) {
		return &v
	}

	return nil
}

var _ error = (*ValidationError)(nil)
