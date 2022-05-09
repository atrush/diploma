package psql

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"path"
	"path/filepath"
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

// // Set migrations to db
// func initMigrations(db *sql.DB, migrationsPath string) error {
// 	driver, err := postgres.WithInstance(db, &postgres.Config{})
// 	if err != nil {
// 		return fmt.Errorf("ошибка миграции бд:%w", err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
// 	if err != nil {
// 		return fmt.Errorf("ошибка миграции бд:%w", err)
// 	}

// 	if err = m.Up(); err != nil {
// 		return fmt.Errorf("ошибка миграции бд:%w", err)
// 	}

// 	return nil
// }

// getFixturesDir returns current file directory.
func getMigrationsFolder() string {
	_, filePath, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return path.Dir(filePath)
}

func getMigrationsRelPath() string {

	p, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	//dir := filepath.Dir(migrationFile())
	//cc := os.DirFS(dir)
	//return cc.Open()
	dir := getMigrationsFolder()
	//if vol := filepath.VolumeName(dir); vol != "" {
	//	root = vol
	//}
	rel, err := filepath.Rel(p, dir)
	if err != nil {
	}
	rel = "file://" + filepath.ToSlash(rel)
	//dd, err := cdup.FindIn(os.DirFS(root), rel, ".git")
	//log.Println(dd)
	return rel
}
