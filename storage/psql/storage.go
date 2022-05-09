package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/atrush/diploma.git/storage"
	"github.com/atrush/diploma.git/storage/psql/fixtures"
	"github.com/atrush/diploma.git/storage/psql/migrations"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	orderRepo    *orderRepository
	userRepo     *userRepository
	db           *sql.DB
	conStringDSN string
}

//  NewStorage inits new connection to psql storage.
//  !!!! On init drop all and init tables.
func NewStorage(dsn string) (*Storage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("ошибка инициализации бд:%v", "строка соединения с бд пуста")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	//if err := initBase(db); err != nil {
	//	return nil, err
	//}
	if err := migrations.RunMigrations(db, "tstdb"); err != nil {
		return nil, err
	}

	st := &Storage{
		db:           db,
		conStringDSN: dsn,
	}

	st.orderRepo = newOrderRepository(db)
	st.userRepo = newUserRepository(db)

	return st, nil
}

func NewTestStorage(dsn string) (*Storage, error) {

	st, err := NewStorage(dsn)
	if err != nil {
		return nil, err
	}

	if err := fixtures.LoadFixtures(context.Background(), st.db); err != nil {
		return nil, err
	}

	return st, nil
}

//  User returns users repository.
func (s *Storage) User() storage.UserRepository {
	return s.userRepo
}

//  Order returns users repository.
func (s *Storage) Order() storage.OrderRepository {
	return s.orderRepo
}

//  Close  closes database connection.
func (s Storage) Close() {
	if s.db == nil {
		return
	}

	s.db.Close()
	s.db = nil
}
