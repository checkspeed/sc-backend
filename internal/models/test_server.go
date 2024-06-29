package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateTestServer struct {
	ID         string `json:"-"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

type TestServer struct {
	ID         string         `json:"id"`
	Identifier string         `json:"identifier"`
	Name       string         `json:"name"`
	City       string         `json:"city"`
	Country    string         `json:"country"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}
