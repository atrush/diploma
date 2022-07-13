package handler

import (
	"fmt"
	"github.com/atrush/diploma.git/api/model"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/atrush/diploma.git/services/order"
	"github.com/atrush/diploma.git/services/withdraw"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	svcAuth     auth.Authenticator
	svcOrder    order.OrderManager
	svcWithdraw withdraw.WithdrawManager
}

// NewHandler Return new handler
func NewHandler(auth auth.Authenticator, order order.OrderManager, withdraw withdraw.WithdrawManager) (*Handler, error) {

	return &Handler{
		svcAuth:     auth,
		svcOrder:    order,
		svcWithdraw: withdraw,
	}, nil
}

func (h *Handler) GetUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	ctxID := r.Context().Value(model.ContextKeyUserID)
	if ctxID == nil {
		return uuid.Nil, fmt.Errorf("ail to get user id from context: user id is empty")
	}

	userID, err := uuid.Parse(ctxID.(string))
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to get user id from context: %w", err)
	}

	return userID, nil
}
