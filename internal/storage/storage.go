package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(connString string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}
	return &Storage{conn: conn}, nil
}

func (s *Storage) Conn() *pgx.Conn {
	return s.conn
}


func (s *Storage) Query(query string, args ...interface{}) (pgx.Rows, error) {
	return s.conn.Query(context.Background(), query, args...)
}

func (s *Storage) Exec(query string, args ...interface{}) (pgconn.CommandTag, error) {
	return s.conn.Exec(context.Background(), query, args...)
}

