package handler

import (
	"encoding/json"
	"errors"
	apimodel "github.com/atrush/diploma.git/api/model"
	"github.com/atrush/diploma.git/model"
	"io/ioutil"
	"net/http"
)

//  GetBalance returns user balance,
//  200 — success,
//  401 — user not authenticated,
//  500 — server error.
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// context must contain user id, if not its internal error
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := h.svcWithdraw.GetBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsResult, err := json.Marshal(apimodel.BalanceResponseFromCanonical(balance))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(jsResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//  WithdrawAddToUser adds withdraw to user.
//  200 — success;
//  401 — user not authenticated,
//  400 — wrong request,
//  402 — low founds,
//  422 — wrong number
//  500 — server error.
func (h *Handler) WithdrawAddToUser(w http.ResponseWriter, r *http.Request) {
	// context must contain user id, if not its internal error
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// read withdraw from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rWithdraw apimodel.WithdrawRequest
	if err := json.Unmarshal(body, &rWithdraw); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//  422 if not valid luhn
	if !rWithdraw.NumberIsValidLuhn() {
		http.Error(w, "bad request", http.StatusUnprocessableEntity)
		return
	}

	withdraw := rWithdraw.ToCanonical(userID)

	_, err = h.svcWithdraw.Create(r.Context(), withdraw)
	if err != nil && !errors.Is(err, model.ErrorWithdrawExist) {
		//  402 if low founds
		if errors.Is(err, model.ErrorNotEnoughFounds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}

		//  500 is internal error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  202 if order saved
	w.WriteHeader(http.StatusOK)
}

//  WithdrawsGetListForUser returns user withdraws list
//  200 — if orders exist, returns list
//  204 — no withdraws for user
//  401 — user not authenticated,
//  500 — server error.
func (h *Handler) WithdrawsGetListForUser(w http.ResponseWriter, r *http.Request) {
	// context must contain user id, if not its internal error
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	withdraws, err := h.svcWithdraw.GetForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdraws) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	jsResult, err := json.Marshal(apimodel.WithdrawResponseListFromCanonical(withdraws))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(jsResult); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
