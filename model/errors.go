package model

import "errors"

var (
	ErrorOrderExistAnotheUser = errors.New("order already exist for another user")
	ErrorOrderExist           = errors.New("order already exist for that user")
	ErrorConflictSaveUser     = errors.New("user already exist")
	ErrorItemNotFound         = errors.New("item not found")
)
