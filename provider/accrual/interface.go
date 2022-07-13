package accrual

import (
	"context"
	"github.com/atrush/diploma.git/model"
)

//  AccrualProvider loads Accrual for order from bonus subsystem
type AccrualProvider interface {
	Get(ctx context.Context, number string) (model.Accrual, error)
}
