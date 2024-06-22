package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrator struct {
	db *sql.DB
}

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

func NewMigrator(store Store) (*migrator, error) {
	sqlDb, err := store.DB().DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql db: %v", err)
	}

	return &migrator{sqlDb}, nil
}

func (m *migrator) Up(ctx context.Context) error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}
	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = mg.Up()
	if err != nil {
		return fmt.Errorf("failed to run up migration: %v", err)
	}
	fmt.Println("database migrated")
	return nil
}

func (m *migrator) Down(ctx context.Context) error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}
	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = mg.Down()
	if err != nil {
		return fmt.Errorf("failed to run down migration: %v", err)
	}
	fmt.Println("database migrated")
	return nil
}
