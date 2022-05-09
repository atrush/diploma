package storage

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/google/uuid"
)

//  Storage is the interface that wraps methods for working with the database.
type Storage interface {

	//  User returns repository for working with users.
	User() UserRepository
	//  Order returns repository for working with orders.
	Order() OrderRepository
	//  Close closes storage connection.
	Close()
}

type UserRepository interface {
	//  Adds new user to storage
	Create(ctx context.Context, user model.User) (model.User, error)
	//  Returns user from storage
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type OrderRepository interface {
	//  Create adds new order to database? if not exist.
	//  If exist with that number for user returns ErrorConflictSaveOrder.
	Create(ctx context.Context, order model.Order) (model.Order, error)
	//  GetForUser selects user orders
	GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error)
	//  UpdateStatus updates status for order record, selected by id.
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
	//  UpdateAccrual updates order accrual and status.
	UpdateAccrual(ctx context.Context, id uuid.UUID, status model.OrderStatus, accrual int) error
	//  GetUnprocessedOrders gets unprocessed orders, and change status to UPDATING
	GetUnprocessedOrders(ctx context.Context, limit int) ([]model.Order, error)
	//  UpdateStatusToNewBatch Updates statuse to new and sets accrual to 0, for batch orders
	UpdateStatusToNewBatch(ctx context.Context, batch []model.Order) (err error)
}
