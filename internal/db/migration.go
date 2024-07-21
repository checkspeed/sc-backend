package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrator struct {
	db *sql.DB
}

type Migrator interface {
	Up(ctx context.Context, migrationsPaths ...string) error
	Down(ctx context.Context, migrationsPaths ...string) error
}

func NewMigrator(store Store) (*migrator, error) {
	sqlDb, err := store.DB().DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve sql db: %v", err)
	}

	return &migrator{sqlDb}, nil
}

func (m *migrator) Up(ctx context.Context, migrationsPaths ...string) error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}

	migrationsURL, err := m.constructMigrationPath(migrationsPaths...)
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = mg.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run up migration: %v", err)
	}

	fmt.Println("database migrated up")
	return nil
}

func (m *migrator) Down(ctx context.Context, migrationsPaths ...string) error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}

	migrationsURL, err := m.constructMigrationPath(migrationsPaths...)
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = mg.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run down migration: %v", err)
	}

	fmt.Println("database migrated down")
	return nil
}

func (m *migrator) constructMigrationPath(migrationsPaths ...string) (string, error) {
	var migrationsPath string

	if len(migrationsPaths) > 0 && migrationsPaths[0] != "" {
		migrationsPath = migrationsPaths[0]
	}

	fmt.Println("mig path", migrationsPath)
	if migrationsPath == "" {
		_, testFilePath, _, _ := runtime.Caller(0)
		testDir := filepath.Dir(testFilePath)

		// Default to the migrations folder at the same level as the executable
		migrationsPath = filepath.Join(testDir, "../db/migrations")
	}

	// Convert to file URL format
	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)
	if strings.HasPrefix(migrationsURL, "/") {
		migrationsURL = "file://" + migrationsURL
	}

	return migrationsURL, nil
}
