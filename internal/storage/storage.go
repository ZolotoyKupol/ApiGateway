package storage

import (
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(connString string) (*Storage, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return &Storage{db: db}, nil
}

func NewStorageWithDB(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}

func (s *Storage) Close() error {
	db, err := s.db.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get database connection")
	}
	return db.Close()
}
