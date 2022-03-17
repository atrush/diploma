package api

import (
	"github.com/atrush/diploma.git/services/auth"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

type Handler struct {
	svcAuth   auth.Authenticator
	tokenAuth *jwtauth.JWTAuth
}

var ContextKeyUserID = contextKey("user-id")

// NewHandler Return new handler
func NewHandler(auth auth.Authenticator, t *jwtauth.JWTAuth) (*Handler, error) {

	return &Handler{
		svcAuth:   auth,
		tokenAuth: t,
	}, nil
}

// Ok return ok status
func (h *Handler) Ok(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
