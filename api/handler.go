package api

import (
	"github.com/atrush/diploma.git/services/auth"
	"net/http"
)

type Handler struct {
	svcAuth auth.Authenticator
}

var ContextKeyUserID = contextKey("user-id")

// NewHandler Return new handler
func NewHandler(auth auth.Authenticator) (*Handler, error) {

	return &Handler{
		svcAuth: auth,
	}, nil
}

// Ok return ok status
func (h *Handler) Ok(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
