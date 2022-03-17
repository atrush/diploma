package storage

import "errors"

var (
	ErrorConflictSaveUser = errors.New("user already exist")
	ErrorItemNotFound     = errors.New("item not found")
)
