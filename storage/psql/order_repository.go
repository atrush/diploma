package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
	"log"
)

var _ storage.OrderRepository = (*orderRepository)(nil)

//  orderRepository implements OrderRepository interface, provides actions with order records in psql storage.
type orderRepository struct {
	db *sql.DB
}

//  newOrderRepository inits new order repository.
func newOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{
		db: db,
	}
}

// UpdateStatus implements interface form implements OrderRepository
func (o *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	if _, err := o.db.ExecContext(
		ctx,
		"UPDATE orders SET status = $1 WHERE id = $2",
		status, id); err != nil {
		return err
	}

	return nil
}

// UpdateAccrual implements interface form implements OrderRepository
func (o *orderRepository) UpdateAccrual(ctx context.Context, id uuid.UUID, status model.OrderStatus, accrual int) error {
	if _, err := o.db.ExecContext(
		ctx,
		"UPDATE orders SET status = $1, accrual= $2  WHERE id = $3",
		status, accrual, id); err != nil {
		return err
	}

	return nil
}

//  Create implements OrderRepository Create interface
func (o *orderRepository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	//  init transaction
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Order{}, err
	}

	defer func() {
		tx.Rollback()
		//todo: logging if err
	}()

	//  check order exist for user
	userID := uuid.Nil
	err = tx.QueryRowContext(
		ctx,
		`SELECT user_id FROM orders WHERE number = $1 LIMIT 1`,
		order.Number).Scan(&userID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, err
	}

	if userID != uuid.Nil {
		if userID == order.UserID {
			return model.Order{}, model.ErrorOrderExist
		}
		return model.Order{}, model.ErrorOrderExistAnotherUser
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
func (o *orderRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Order, error) {
	userOrders := make([]model.Order, 0)

	rows, err := o.db.QueryContext(ctx,
		"SELECT id, user_id, number, uploaded_at, status, accrual FROM orders WHERE user_id = $1",
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

func (o *orderRepository) UpdateStatusToNewBatch(ctx context.Context, batch []model.Order) (err error) {
	if len(batch) == 0 {
		return nil
	}

	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if rollErr := tx.Rollback(); rollErr != nil {
			log.Println("error rollback for orders batch status update")
		}
	}()

	for _, o := range batch {
		if _, err := tx.ExecContext(
			ctx,
			"UPDATE orders SET status = $1, accrual= $2  WHERE id = $3",
			model.OrderStatusNew, 0, o.ID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) GetUnprocessedOrders(ctx context.Context, limit int) ([]model.Order, error) {
	userOrders := make([]model.Order, 0)

	//  init transaction
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		tx.Rollback()
		//todo: logging if err
	}()

	rows, err := tx.QueryContext(ctx,
		"UPDATE orders SET status = $1"+
			"WHERE id IN ( SELECT id FROM orders WHERE status = $2 or status = $3 LIMIT $4 ) RETURNING * ",
		model.OrderStatusUpdating,
		model.OrderStatusProcessing,
		model.OrderStatusNew,
		limit,
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

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return userOrders, nil
}

func (o *orderRepository) GetUserAccrualsSum(ctx context.Context, userID uuid.UUID) (int, error) {
	rows, err := o.db.QueryContext(ctx,
		"SELECT accrual FROM orders WHERE user_id = $1 AND status = $2",
		userID,
		model.OrderStatusProcessed,
	)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	sum := 0
	for rows.Next() {
		var w int
		if err = rows.Scan(
			&w,
		); err != nil {
			return 0, err
		}
		sum -= w
	}

	return sum, nil
}
