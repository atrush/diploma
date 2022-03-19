package psql

import (
	"database/sql"
	"github.com/atrush/diploma.git/storage"
	"github.com/atrush/diploma.git/storage/psql/migrations"

	"fmt"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	//shortURLRepo *shortURLRepository
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

//  Close  closes database connection.
func (s Storage) Close() {
	if s.db == nil {
		return
	}

	s.db.Close()
	s.db = nil
}

//  initBase drops all and inits database tables.
func initBase(db *sql.DB) error {
	row := db.QueryRow("DROP SCHEMA public CASCADE;CREATE SCHEMA public;")
	if row.Err() != nil {
		return row.Err()
	}
	_, err := db.Exec("create extension if not exists \"uuid-ossp\";" +
		"CREATE TABLE IF NOT EXISTS users (" +
		"		id uuid primary key default uuid_generate_v4()," +
		"		login varchar (255) unique not null," +
		"		passhash varchar (60) not null);")
	if row.Err() != nil {
		return err
	}
	return nil

}
