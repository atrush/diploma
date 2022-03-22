package model

import (
	"github.com/google/uuid"
	"time"
)

type (
	OrderStatus   string
	AccrualStatus string
)

//  result money accuracy coeff
const MoneyAccuracy = 1000

var (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusRegistered OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"

	AccrualStatusRegistered AccrualStatus = "REGISTERED"
	AccrualStatusInvalid    AccrualStatus = "INVALID"
	AccrualStatusProcessing AccrualStatus = "PROCESSING"
	AccrualStatusProcessed  AccrualStatus = "PROCESSED"
)

func (a AccrualStatus) IsValid() bool {
	statuses := map[string]struct{}{
		"REGISTERED": struct{}{},
		"INVALID":    struct{}{},
		"PROCESSING": struct{}{},
		"PROCESSED":  struct{}{},
	}
	_, ok := statuses[string(a)]
	return ok
}

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

type Accrual struct {
	Number  string
	Status  AccrualStatus
	Accrual int
}
