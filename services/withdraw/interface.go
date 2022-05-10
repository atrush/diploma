package withdraw

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/google/uuid"
)

type WithdrawManager interface {
	//  Create adds new withdraw to user
	//  checks is enogh founds
	Create(ctx context.Context, withdraw model.Withdraw) (model.Withdraw, error)
	//  GetForUser returns withdraws for user
	GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Withdraw, error)
}
