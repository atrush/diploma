package order

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/google/uuid"
)

//  OrderManager is the interface that wraps methods for order processing.
type OrderManager interface {
	//  AddToUser adds new order to user
	AddToUser(ctx context.Context, number string, userID uuid.UUID) (model.Order, error)
	//  GetForUser gets user orders
	GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	//  UpdateStatus updates order status
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
	//  UpdateStatus updates order accrual and status
	UpdateAccrual(ctx context.Context, order model.Order, accrual model.Accrual) error
}
