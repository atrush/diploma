package handler

import (
	apimiddleware "github.com/atrush/diploma.git/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter(handler *Handler) *chi.Mux {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/user/register", handler.Register)
		r.Post("/api/user/login", handler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(apimiddleware.MiddlewareAuth)

		//  orders
		r.Post("/api/user/orders", handler.OrderAddToUser)
		r.Get("/api/user/orders", handler.OrderGetListForUser)

		//  withdrawals
		r.Post("/api/user/balance/withdraw", handler.WithdrawAddToUser)
		r.Get("/api/user/balance/withdrawals", handler.WithdrawsGetListForUser)

		//  balance
		r.Get("/api/user/balance", handler.GetBalance)
	})

	return r
}
