package storage

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")

	// News
	ErrNewsAlreadyExists = fmt.Errorf("news %w", ErrAlreadyExists)
	ErrNewsNotFound      = fmt.Errorf("news %w", ErrNotFound)
)
