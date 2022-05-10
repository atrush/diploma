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
		// POST /api/user/register — регистрация пользователя;
		r.Post("/api/user/register", handler.Register)
		// POST /api/user/login — аутентификация пользователя;
		r.Post("/api/user/login", handler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(apimiddleware.MiddlewareAuth)
		r.Post("/api/user/orders", handler.OrderAddToUser)
		r.Get("/api/user/orders", handler.OrderGetListForUser)
		r.Post("/api/user/balance/withdraw", handler.WithdrawAddToUser)
		r.Get("/api/user/balance/withdrawals", handler.WithdrawsGetListForUser)
	})
	//// POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	//r.Post("/api/user/orders", handler.Ok)
	//// GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	//r.Get("/api/user/orders", handler.Ok)
	//// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
	//r.Get("/api/user/balance", handler.Ok)
	//// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	//r.Post("/api/user/withdraw", handler.Ok)
	//// GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	//r.Get("/api/user/withdrawals", handler.Ok)
	return r
}
