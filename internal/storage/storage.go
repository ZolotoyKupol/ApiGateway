package storage

import (
	"apigateway/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(connString string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.Guest{}, &models.Room{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}