package psql

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"path"
	"runtime"
)

//  RunMigrations apply migrations to database
//  !! drop database before run migrations
func RunMigrations(dbDSN string, dbName string) error {
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		return err
	}

	//m.Drop()

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}

// getFixturesDir returns current file directory.
func getMigrationsFolder() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return path.Dir(filePath)
}
