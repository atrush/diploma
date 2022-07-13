package handler

import (
	"encoding/json"
	"errors"
	apimodel "github.com/atrush/diploma.git/api/model"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/pkg/validation"
	"io/ioutil"
	"net/http"
)

//  OrderAddToUser adds order to user.
//  200 — number of order is exist for user,
//  202 — number of order accepted,
//  400 — wrong request,
//  401 — user not authenticated,
//  409 — number of order was uploaded by another user,
//  422 — wrong order number,
//  500 — server error.
func (h *Handler) OrderAddToUser(w http.ResponseWriter, r *http.Request) {
	// context must contain user id, if not its internal error
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// read number
	number, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	strNumber := string(number)

	//  400 if empty number or no digits
	if len(number) == 0 || !validation.IsInt(strNumber) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//  422 if not valid luhn
	if !validation.ValidLuhn(strNumber) {
		http.Error(w, "bad request", http.StatusUnprocessableEntity)
		return
	}

	//  create new order
	_, err = h.svcOrder.AddToUser(r.Context(), strNumber, userID)
	if err != nil {
		//  409 if exist for another user
		if errors.Is(err, model.ErrorOrderExistAnotherUser) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		//  200 if exist for current user
		if errors.Is(err, model.ErrorOrderExist) {
			http.Error(w, err.Error(), http.StatusOK)
			return
		}
		//  500 if something wrong
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  202 if order saved
	w.WriteHeader(http.StatusAccepted)
}

//  OrderGetListForUser returns user orders list
//  200 — if orders exist, returns list
//  204 — no orders for user
//  401 — user not authenticated,
//  500 — server error.
func (h *Handler) OrderGetListForUser(w http.ResponseWriter, r *http.Request) {
	// context must contain user id, if not its internal error
	userID, err := h.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orders, err := h.svcOrder.GetForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	jsResult, err := json.Marshal(apimodel.OrderResponseListFromCanonical(orders))
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
