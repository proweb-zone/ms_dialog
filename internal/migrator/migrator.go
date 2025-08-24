package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const migrationsDir = "migrations"

type Migrator struct {
	srcDriver    source.Driver
	databaseName string
}

func MustGetNewMigrator(databaseName string) *Migrator {
	d, err := iofs.New(MigrationsFS, migrationsDir)
	if err != nil {
		panic(err)
	}
	return &Migrator{
		srcDriver:    d,
		databaseName: databaseName,
	}
}

func (m *Migrator) Up(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, m.databaseName, driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	defer func() {
		migrator.Close()
	}()

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("no new migrations to apply: %v \n", migrate.ErrNoChange)
			return nil
		}
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}

func (m *Migrator) Down(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, m.databaseName, driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	defer func() {
		migrator.Close()
	}()

	if err = migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("no new migrations to apply: %v \n", migrate.ErrNoChange)
			return nil
		}
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}

func (m *Migrator) UpAndKeepConnection(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, m.databaseName, driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("no new migrations to apply: %v \n", migrate.ErrNoChange)
			return nil
		}
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}
