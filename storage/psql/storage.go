package psql

import (
	"database/sql"
	"github.com/atrush/diploma.git/storage"
	"github.com/atrush/diploma.git/storage/psql/migrations"

	"fmt"
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
func NewStorage(conStringDSN string) (*Storage, error) {
	if conStringDSN == "" {
		return nil, fmt.Errorf("ошибка инициализации бд:%v", "строка соединения с бд пуста")
	}

	db, err := sql.Open("postgres", conStringDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	//if err := initBase(db); err != nil {
	//	return nil, err
	//}
	if err := migrations.RunMigrations(db); err != nil {
		return nil, err
	}

	st := &Storage{
		db:           db,
		conStringDSN: conStringDSN,
	}

	//st.shortURLRepo = newShortURLRepository(db)
	st.userRepo = newUserRepository(db)

	return st, nil
}

//  User returns users repository.
func (s *Storage) User() storage.UserRepository {
	return s.userRepo
}

//  Order returns users repository.
func (s *Storage) Order() storage.OrderRepository {
	return s.Order()
}

//  Close  closes database connection.
func (s Storage) Close() {
	if s.db == nil {
		return
	}

	s.db.Close()
	s.db = nil
}
