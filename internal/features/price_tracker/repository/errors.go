package repository

import "errors"

var (
	ErrUrlExists = errors.New("Item with this URL already exists")
	ErrItemNotFound = errors.New("Item not found")
)