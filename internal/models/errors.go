package models

import "errors"

var (
	ErrNoRecord       = errors.New("models: no matching record found")
	ErrDuplicateEntry = errors.New("models: duplicate entry")
)
