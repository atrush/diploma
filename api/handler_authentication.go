package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/google/uuid"
	"net/http"
)

//  Register registers user, sets cookie with jwt token.
//  200 — user registered;
//  400 — wrong request format;
//  409 — user exist;
//  500 — internal server error.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	loginData, err := h.readLoginRequest(w, r)
	if err != nil {
		return
	}

	//  authenticate and get user.
	user, err := h.svcAuth.CreateUser(r.Context(), loginData.Login, loginData.Password)
	if err != nil {
		//  if exist returns 409.
		if errors.Is(err, auth.ErrorUserAlreadyExist) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		//  if something was wrong 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  set token to response
	if err := h.setToken(w, user.ID); err != nil {
		return
	}

}

//  Login authenticates user, sets jwt token.
//  200 — user authenticated;
//  400 — wrong request format;
//  401 — wrong login/password;
//  500 — internal server error.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	loginData, err := h.readLoginRequest(w, r)
	if err != nil {
		return
	}

	//  authenticate and get user.
	user, err := h.svcAuth.Authenticate(r.Context(), loginData.Login, loginData.Password)
	if err != nil {
		//  if wrong login,password return 401.
		if errors.Is(err, auth.ErrorWrongAuthData) {
			http.Error(w, "incorrect login or password", http.StatusUnauthorized)
			return
		}

		//  if something was wrong 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  set token to response
	if err := h.setToken(w, user.ID); err != nil {
		return
	}

}

//  readLoginRequest reads login data from request.
func (h Handler) readLoginRequest(w http.ResponseWriter, r *http.Request) (LoginRequest, error) {
	var loginData LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return LoginRequest{}, fmt.Errorf("wrong data format")
	}
	defer r.Body.Close()
	if err := loginData.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return LoginRequest{}, fmt.Errorf("wrong data format")
	}
	return loginData, nil
}

//  setToken sets jwt token with user_id claim to response.
func (h Handler) setToken(w http.ResponseWriter, userID uuid.UUID) error {
	//  encode token with user_id claim.
	token, err := h.svcAuth.EncodeTokenUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return err
	}

	//  set cookie.
	cookie := http.Cookie{
		Name:  "jwt",
		Value: token,
	}
	http.SetCookie(w, &cookie)

	//  set header.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)

	return nil
}
