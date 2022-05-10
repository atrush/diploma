package psql

import (
	"context"
	"database/sql"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
)

var _ storage.WithdrawRepository = (*withdrawRepository)(nil)

//  withdrawRepository implements WithdrawRepository interface, provides actions with order records in psql storage.
type withdrawRepository struct {
	db *sql.DB
}

//  newWithdrawRepository inits new withdraw repository.
func newWithdrawRepository(db *sql.DB) *withdrawRepository {
	return &withdrawRepository{
		db: db,
	}
}

//  Create implements WithdrawRepository Create interface
func (r *withdrawRepository) Create(ctx context.Context, withdraw model.Withdraw) (model.Withdraw, error) {
	//  init transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Withdraw{}, err
	}

	defer func() {
		tx.Rollback()
		//todo: logging if err
	}()

	//  check withdraw exist for user
	c := 0
	err = tx.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM withdraws WHERE number = $2 AND user_id = $1`,
		withdraw.UserID,
		withdraw.Number).Scan(&c)

	if err != nil {
		return model.Withdraw{}, err
	}

	if c > 0 {
		return model.Withdraw{}, model.ErrorWithdrawExist
	}

	sum := 0

	//  calc user accruals
	rowsOrders, err := tx.QueryContext(ctx,
		"SELECT accrual FROM orders WHERE user_id = $1 AND status = $2",
		withdraw.UserID,
		model.OrderStatusProcessed,
	)
	if err != nil {
		return model.Withdraw{}, err
	}

	defer rowsOrders.Close()

	for rowsOrders.Next() {
		var a int
		if err = rowsOrders.Scan(
			&a,
		); err != nil {
			return model.Withdraw{}, err
		}
		sum += a
	}

	//  calc user withdraws
	rowsWithdraws, err := tx.QueryContext(ctx,
		"SELECT sum FROM withdraws WHERE user_id = $1",
		withdraw.UserID,
	)
	if err != nil {
		return model.Withdraw{}, err
	}

	defer rowsWithdraws.Close()

	for rowsWithdraws.Next() {
		var w int
		if err = rowsWithdraws.Scan(
			&w,
		); err != nil {
			return model.Withdraw{}, err
		}
		sum -= w
	}

	if withdraw.Sum > sum {
		return model.Withdraw{}, model.ErrorNotEnoughFounds
	}

	// insert withdraw
	if err := tx.QueryRowContext(
		ctx,
		"INSERT INTO withdraws (user_id, number, uploaded_at, sum) VALUES ($1, $2,$3,$4) "+
			"RETURNING id, user_id, number, uploaded_at, sum",
		withdraw.UserID,
		withdraw.Number,
		withdraw.UploadedAt,
		withdraw.Sum,
	).Scan(
		&withdraw.ID,
		&withdraw.UserID,
		&withdraw.Number,
		&withdraw.UploadedAt,
		&withdraw.Sum,
	); err != nil {
		return model.Withdraw{}, err
	}

	//  commit transaction
	if err := tx.Commit(); err != nil {
		return model.Withdraw{}, err
	}

	return withdraw, nil
}

//  GetForUser implements WithdrawRepository GetForUser interface
func (s *withdrawRepository) GetForUser(ctx context.Context, userID uuid.UUID) ([]model.Withdraw, error) {
	userWithdraws := make([]model.Withdraw, 0)

	rows, err := s.db.QueryContext(ctx,
		"SELECT id, user_id, number, uploaded_at, sum FROM withdraws WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var el model.Withdraw
		if err = rows.Scan(
			&el.ID,
			&el.UserID,
			&el.Number,
			&el.UploadedAt,
			&el.Sum,
		); err != nil {
			return nil, err
		}

		userWithdraws = append(userWithdraws, el)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userWithdraws, nil
}
