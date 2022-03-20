package psql

import (
	"context"
	"database/sql"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
)

var _ storage.OrderRepository = (*orderRepository)(nil)

//  orderRepository implements OrderRepository interface, provides actions with order records in psql storage.
type orderRepository struct {
	db *sql.DB
}

//  newOrderRepository inits new order repository.
func newOrderRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

//  Create implements OrderRepository Create interface
func (r *orderRepository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	//  init transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Order{}, err
	}

	defer func() {
		tx.Rollback()
		//todo: logging if err
	}()

	//  check exist for user
	userID := uuid.Nil
	if err := tx.QueryRowContext(
		ctx,
		`SELECT id FROM orders WHERE number = $1 AND user_id = $2 LIMIT 1`,
		order.UserID,
		order.Number).Scan(&userID); err != nil {
		return model.Order{}, err
	}

	if userID != uuid.Nil {
		if userID == order.UserID {
			return model.Order{}, model.ErrorOrderExist
		}
		return model.Order{}, model.ErrorOrderExistAnotheUser
	}

	// insert
	if err := tx.QueryRowContext(
		ctx,
		"INSERT INTO orders (user_id, number, uploaded_at, status, accrual) "+
			"VALUES ($1, $2,$3,$4,$5) "+
			"RETURNING id, user_id, number, uploaded_at, status, accrual",
		order.UserID,
		order.Number,
		order.UploadedAt,
		order.Status,
		order.Accrual,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.UploadedAt,
		&order.Status,
		&order.Accrual,
	); err != nil {
		return model.Order{}, err
	}

	//  commit transaction
	if err := tx.Commit(); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

//  GetForUser implements OrderRepository GetForUser interface
func (s *orderRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	userOrders := make([]model.Order, 0)

	rows, err := s.db.QueryContext(ctx,
		`SELECT  id, user_id, number, uploaded_at, status, accrual
		FROM orders WHERE user_id = $1
        ORDER BY uploaded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order model.Order
		if err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.UploadedAt,
			&order.Status,
			&order.Accrual,
		); err != nil {
			return nil, err
		}

		userOrders = append(userOrders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userOrders, nil
}
