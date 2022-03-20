package order

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
	"time"
)

var _OrderManager = (*Order)(nil)

type Order struct {
	storage storage.Storage
}

func NewOrder(s storage.Storage) (*Order, error) {
	return &Order{
		storage: s,
	}, nil
}

func (o *Order) AddToUser(ctx context.Context, number string, userID uuid.UUID) (model.Order, error) {
	order := model.Order{
		UserID:     userID,
		Number:     number,
		Accrual:    0,
		Status:     model.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	order, err := o.storage.Order().Create(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *Order) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	userOrders, err := o.storage.Order().GetForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}
