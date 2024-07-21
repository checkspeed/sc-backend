package db

import (
	"context"
	"fmt"

	_ "embed"

	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/database/postgres"
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

func NewStore(dbUrl string) (*store, error) {
	db, err := gorm.Open(gPostgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return &store{}, err
	}

	fmt.Println("database connected")

	return &store{db}, nil
}

func (s *store) DB() *gorm.DB {
	return s.db
}

func (s store) CloseConn(ctx context.Context) error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
