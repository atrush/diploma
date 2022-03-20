package api

import (
	"fmt"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/atrush/diploma.git/services/order"
	"github.com/google/uuid"
	"net/http"
)

type Handler struct {
	svcAuth  auth.Authenticator
	svcOrder order.OrderManager
}

var ContextKeyUserID = contextKey("user-id")

// NewHandler Return new handler
func NewHandler(auth auth.Authenticator, order order.OrderManager) (*Handler, error) {

	return &Handler{
		svcAuth:  auth,
		svcOrder: order,
	}, nil
}

// Ok return ok status
func (h *Handler) Ok(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	ctxID := r.Context().Value(ContextKeyUserID)
	if ctxID == nil {
		return uuid.Nil, fmt.Errorf("ail to get user id from context: user id is empty")
	}

	userID, err := uuid.Parse(ctxID.(string))
	if err != nil {
		return uuid.Nil, fmt.Errorf("fail to get user id from context: %w", err)
	}

	return userID, nil
}
