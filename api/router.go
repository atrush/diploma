package api

import (
	"github.com/atrush/diploma.git/api/handler"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *handler.Handler) *chi.Mux {
	r := chi.NewRouter()
	// POST /api/user/register — регистрация пользователя;
	r.Post("/api/user/register", handler.Ok)
	// POST /api/user/login — аутентификация пользователя;
	r.Post("/api/user/login", handler.Ok)

	// POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
	r.Post("/api/user/orders", handler.Ok)
	// GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	r.Get("/api/user/orders", handler.Ok)
	// GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
	r.Get("/api/user/balance", handler.Ok)
	// POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	r.Post("/api/user/withdraw", handler.Ok)
	// GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	r.Get("/api/user/withdrawals", handler.Ok)
	return r
}
