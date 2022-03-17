package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		// POST /api/user/register — регистрация пользователя;
		r.Post("/api/user/register", handler.Register)
		// POST /api/user/login — аутентификация пользователя;
		r.Post("/api/user/login", handler.Login)
	})

	//r.Group(func(r chi.Router) {
	//	r.Use(jwtauth.Verifier(handler.tokenAuth))
	//	r.Get("/api/user/check", handler.GetUserIDFromContext)
	//})
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
