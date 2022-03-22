package order

import (
	"context"
	"fmt"
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

func (o *Order) UpdateAccrual(ctx context.Context, order model.Order, accrual model.Accrual) error {
	if accrual.Status == model.AccrualStatusRegistered {
		order.Status = model.OrderStatusProcessing
	}
	if accrual.Status == model.AccrualStatusInvalid {
		order.Status = model.OrderStatusInvalid
	}
	if accrual.Status == model.AccrualStatusProcessing {
		order.Status = model.OrderStatusProcessing
	}
	if accrual.Status == model.AccrualStatusProcessed {
		order.Status = model.OrderStatusProcessed
		order.Accrual = accrual.Accrual
	}

	if err := o.storage.Order().UpdateAccrual(ctx, order.ID, order.Status, order.Accrual); err != nil {
		return fmt.Errorf("error updating uccrual:%w", err)
	}

	return nil
}

func (o *Order) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	return o.storage.Order().UpdateStatus(ctx, id, status)
}

func (o *Order) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	userOrders, err := o.storage.Order().GetForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}
