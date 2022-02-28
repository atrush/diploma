package api

import "net/http"

type Handler struct {
}

// NewHandler Return new handler
func NewHandler() (*Handler, error) {
	return &Handler{}, nil
}

// Ok return ok status
func (h *Handler) Ok(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
