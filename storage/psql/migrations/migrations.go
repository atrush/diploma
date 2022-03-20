package migrations

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

//  RunMigrations apply migrations to database
//  !! drop database before run migrations
func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://../../storage/psql/migrations", "praktikum", driver)
	if err != nil {
		return err
	}

	if err := m.Drop(); err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
