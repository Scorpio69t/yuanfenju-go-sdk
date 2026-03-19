package yuanfenju

import (
	"errors"
	"fmt"
)

var ErrValidation = errors.New("validation failed")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Field == "" {
		return fmt.Sprintf("%v: %s", ErrValidation, e.Message)
	}
	return fmt.Sprintf("%v: field=%s %s", ErrValidation, e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

func newRequiredFieldError(field string) error {
	return &ValidationError{Field: field, Message: "is required"}
}

func newEnumFieldError(field, value string, allowed []string) error {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("must be one of %v, got %q", allowed, value),
	}
}

func inSet(v string, allowed []string) bool {
	for _, x := range allowed {
		if v == x {
			return true
		}
	}
	return false
}
