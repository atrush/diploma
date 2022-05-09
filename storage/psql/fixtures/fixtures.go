package fixtures

import (
	"context"
	"database/sql"
	"github.com/atrush/diploma.git/model/testdata"
)

// LoadFixtures load fixtures to DB and returns DB objects aggregate.
func LoadFixtures(ctx context.Context, db *sql.DB) error {

	data, err := testdata.ReadTestData()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(
		ctx,
		"DELETE FROM orders; DELETE FROM users; ",
	)
	if err != nil {
		return err
	}

	if len(data.Users) > 0 {
		for _, el := range data.Users {
			_, err := db.ExecContext(
				ctx,
				"INSERT INTO users (id, login, pass_hash) VALUES ($1, $2, $3)",
				el.ID,
				el.Login,
				el.PasswordHash,
			)
			if err != nil {
				return err
			}

		}
	}
	if len(data.Orders) > 0 {
		for _, el := range data.Orders {
			_, err := db.ExecContext(
				ctx,
				"INSERT INTO orders (id,user_id, number, uploaded_at, status, accrual) VALUES ($1, $2,$3,$4,$5,$6) ",
				el.ID,
				el.UserID,
				el.Number,
				el.UploadedAt,
				el.Status,
				el.Accrual,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
