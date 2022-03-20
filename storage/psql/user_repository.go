package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/atrush/diploma.git/model"
	"github.com/atrush/diploma.git/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

var _ storage.UserRepository = (*userRepository)(nil)

//  userRepository implements UserRepository interface, provides actions with user records in psql storage.
type userRepository struct {
	db *sql.DB
}

//  newUserRepository inits new user repository.
func newUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

//  Save saves user to database.
//  If login exist return ErrorConflictSaveUser
func (r *userRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (login, pass_hash) VALUES ($1, $2) RETURNING id, login, passhash",
		user.Login,
		user.PasswordHash,
	).Scan(&user.ID, &user.Login, &user.PasswordHash)

	if err != nil {
		//  if exist return ErrorConflictSaveUser
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == pgerrcode.UniqueViolation && pqErr.Constraint == "users_login_key" {
			return model.User{}, storage.ErrorConflictSaveUser
		}

		return model.User{}, err
	}

	return user, nil
}

//  GetByLogin selects user by login
//	if not found, returns ErrorItemNotFound
func (s *userRepository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	if err := s.db.QueryRowContext(ctx,
		`SELECT id, login, pass_hash FROM users WHERE login = $1`,
		login,
	).Scan(&user.ID, &user.Login, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, storage.ErrorItemNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

//  Exist checks that user is exist in database.
func (r *userRepository) Exist(userID uuid.UUID) (bool, error) {
	count := 0
	err := r.db.QueryRow(
		"SELECT  COUNT(*) as count FROM users WHERE id = $1", userID).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
