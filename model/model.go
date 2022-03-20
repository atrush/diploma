package model

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

var (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
}

type Order struct {
	ID         uuid.UUID
	Number     string
	UserID     uuid.UUID
	Status     OrderStatus
	Accrual    int
	UploadedAt time.Time
}
