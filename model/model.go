package model

import (
	"github.com/google/uuid"
)

type OrderStatus string

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
}
