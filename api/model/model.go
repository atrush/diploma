package model

import (
	"fmt"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/pkg/validation"
	"github.com/google/uuid"
	"time"
)

type (
	LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	OrderResponse struct {
		Number     string  `json:"number"`
		Status     string  `json:"status"`
		Accrual    float64 `json:"accrual,omitempty"`
		UploadedAt string  `json:"uploaded_at"`
	}

	AccrualResponse struct {
		Number  string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual,omitempty"`
	}

	WithdrawRequest struct {
		Number string  `json:"order"`
		Sum    float64 `json:"sum"`
	}

	WithdrawResponse struct {
		Number      string  `json:"order"`
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	BalanceResponse struct {
		Accruals    float64 `json:"current"`
		Withdrawals float64 `json:"withdrawn"`
	}

	ContextKey string
)

var ContextKeyUserID = ContextKey("user-id")

func BalanceResponseFromCanonical(b model.Balance) BalanceResponse {
	return BalanceResponse{
		Accruals:    float64(b.Accruals) / float64(model.MoneyAccuracy),
		Withdrawals: float64(b.Withdrawals) / float64(model.MoneyAccuracy),
	}
}

func (w *WithdrawRequest) ToCanonical(userID uuid.UUID) model.Withdraw {
	return model.Withdraw{
		UserID: userID,
		Number: w.Number,
		Sum:    int(w.Sum * model.MoneyAccuracy),
	}
}

func (w *WithdrawRequest) NumberIsValidLuhn() bool {
	return len(w.Number) > 0 && validation.ValidLuhn(w.Number)
}

func (a *AccrualResponse) ToCanonical() (model.Accrual, error) {
	status := model.AccrualStatus(a.Status)
	if !status.IsValid() {
		return model.Accrual{}, fmt.Errorf("error convert AccrualResponse to Accrual: wrong status %v", a.Status)
	}
	return model.Accrual{
		Status:  status,
		Number:  a.Number,
		Accrual: int(a.Accrual * float64(model.MoneyAccuracy)),
	}, nil
}

//  WithdrawResponseListFromCanonical makes list of WithdrawResponse from canonical withdraw.
func WithdrawResponseListFromCanonical(objs []model.Withdraw) []WithdrawResponse {
	responseArr := make([]WithdrawResponse, 0, len(objs))

	for _, el := range objs {
		o := WithdrawResponse{
			Number:      el.Number,
			Sum:         float64(el.Sum) / float64(model.MoneyAccuracy),
			ProcessedAt: el.UploadedAt.Format(time.RFC3339),
		}

		responseArr = append(responseArr, o)
	}

	return responseArr
}

//  OrderResponseListFromCanonical makes list of OrderResponse from canonical orders.
func OrderResponseListFromCanonical(objs []model.Order) []OrderResponse {
	responseArr := make([]OrderResponse, 0, len(objs))

	for _, order := range objs {
		o := OrderResponse{
			Number:     order.Number,
			Status:     string(order.Status),
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}
		if order.Accrual > 0 {
			o.Accrual = float64(order.Accrual) / float64(model.MoneyAccuracy)
		}

		responseArr = append(responseArr, o)
	}

	return responseArr
}

func (r LoginRequest) Validate() error {
	if len(r.Login) < 3 {
		return fmt.Errorf("login must be larger then 3 symbols")
	}
	if len(r.Login) > 30 {
		return fmt.Errorf("login must be less then 30 symbols")
	}
	if len(r.Password) < 3 {
		return fmt.Errorf("password must be larger then 3 symbols")
	}
	return nil
}
