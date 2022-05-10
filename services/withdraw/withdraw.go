package withdraw

import (
	"context"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
	"time"
)

var _ WithdrawManager = (*Withdraw)(nil)

type Withdraw struct {
	storage storage.Storage
}

func NewWithdraw(s storage.Storage) (*Withdraw, error) {
	return &Withdraw{
		storage: s,
	}, nil
}

func (o *Withdraw) Create(ctx context.Context, withdraw model.Withdraw) (model.Withdraw, error) {
	withdraw.UploadedAt = time.Now()
	return o.storage.Withdraw().Create(ctx, withdraw)
}

func (o *Withdraw) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Withdraw, error) {
	return o.storage.Withdraw().GetForUser(ctx, userID)
}
