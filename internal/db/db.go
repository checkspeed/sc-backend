package db

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

type Store interface {
	CloseConn(ctx context.Context) error
	DB() *gorm.DB
}

func NewStore(dbUrl string) (store, error) {
	db, err := gorm.Open(gPostgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return store{}, err
	}

	fmt.Println("database connected")

	return store{db}, nil
}

func (s *store) RunUpMigration() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve sql db: %v", err)
	}
	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run up migration: %v", err)
	}
	fmt.Println("database migrated")
	return nil
}

func (s *store) RunDownMigration() error {
	sqlDb, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve sql db: %v", err)
	}
	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to retrieve psql instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to retrieve migrator instance: %v", err)
	}
	err = m.Down()
	if err != nil {
		return fmt.Errorf("failed to run down migration: %v", err)
	}
	fmt.Println("database migrated")
	return nil
}

func (s store) CloseConn(ctx context.Context) error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
