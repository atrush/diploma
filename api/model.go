package api

import (
	"fmt"
	"github.com/atrush/diploma.git/model"
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
		Number  string `json:"order"`
		Status  string `json:"status"`
		Accrual int    `json:"accrual,omitempty"`
	}

	contextKey string
)

func (a *AccrualResponse) ToCanonical() (model.Accrual, error) {
	status := model.AccrualStatus(a.Status)
	if !status.IsValid() {
		return model.Accrual{}, fmt.Errorf("error convert AccrualResponse to Accrual: wrong status %v", a.Status)
	}
	return model.Accrual{
		Status:  status,
		Number:  a.Number,
		Accrual: a.Accrual * model.MoneyAccuracy,
	}, nil
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
			o.Accrual = float64(order.Accrual) / model.MoneyAccuracy
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
