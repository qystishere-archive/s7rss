package storage

import (
	"errors"
	"fmt"
	"strings"
)

type CombinedError struct {
	Errors []error
}

func CombineErrors(errs ...error) *CombinedError {
	var ce CombinedError
	for _, err := range errs {
		if err != nil {
			ce.Errors = append(ce.Errors, err)
		}
	}
	if ce.Errors == nil {
		return nil
	}
	return &ce
}

func (ce *CombinedError) Error() string {
	var sb strings.Builder
	for _, err := range ce.Errors {
		if sb.Len() > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")

	// News
	ErrNewsAlreadyExists = fmt.Errorf("news %w", ErrAlreadyExists)
	ErrNewsNotFound      = fmt.Errorf("news %w", ErrNotFound)
)
