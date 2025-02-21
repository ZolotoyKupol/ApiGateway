package repository

import (
	"apigateway/internal/models"
	"apigateway/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type GuestRepo struct {
	store *storage.Storage
	logger *slog.Logger
}

func NewGuestRepo(store *storage.Storage, logger *slog.Logger) *GuestRepo {
	return &GuestRepo{store: store, logger: logger}
}

func (g *GuestRepo) GetGuestsRepo(ctx context.Context) ([]models.Guest, error) {
	rows, err := g.store.Query("SELECT id, first_name, last_name, room_id FROM guests")
	if err != nil {
		g.logger.Error("error fetching guests", "error", err)
		return nil, fmt.Errorf("error fetching guests: %v", err)
	}
	defer rows.Close()

	var guests []models.Guest
	for rows.Next() {
		var guest models.Guest
		err := rows.Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
		if err != nil {
			g.logger.Error("error scanning guest", "error", err)
			return nil, fmt.Errorf("error scanning guest: %v", err)
		}
		guests = append(guests, guest)
	}
	g.logger.Info("successfully fetched guests", "count", len(guests))
	return guests, nil
}

func (g *GuestRepo) CreateGuestRepo(ctx context.Context, guest models.Guest) (string, error) {
	var id string
	err := g.store.Conn().QueryRow(ctx, "INSERT INTO guests (first_name, last_name, room_id) VALUES ($1, $2, $3) RETURNING id", guest.FirstName, guest.LastName, guest.RoomID).Scan(&id)
	if err != nil {
		g.logger.Error("error creating guest", "error", err)
		return "", fmt.Errorf("error creating guest: %v", err)
	}
	g.logger.Info("guest created successfully", "id", id)
	return id, nil
}

func (g *GuestRepo) DeleteGuestRepo(ctx context.Context, id string) (error) {
	_, err := g.store.Conn().Exec(ctx, "DELETE FROM guests WHERE id = $1", id)
	if err != nil {
		g.logger.Error("error deleting guest", err)
		return fmt.Errorf("error deleting guest: %v", err)
	}
	g.logger.Info("guest deleted succesfully", "id", id)
	return nil
}

func (g *GuestRepo) UpdateGuestRepo(ctx context.Context, id string, guest models.Guest) error {
	_, err := g.store.Conn().Exec(ctx, "UPDATE guests SET first_name = $1, last_name = $2, room_id = $3 WHERE id = $4", guest.FirstName, guest.LastName, guest.RoomID, id)
	if err != nil {
		g.logger.Error("error updating guest", "error", err)
		return fmt.Errorf("error updating guest: %v", err)
	}
	g.logger.Info("guest updated successfully", "id", id)
	return nil
}

func (g *GuestRepo) GetGuestByIDRepo(ctx context.Context, id string) (*models.Guest, error) {
	var guest models.Guest
	err := g.store.Conn().QueryRow(ctx, "SELECT id, first_name, last_name, room_id FROM guests WHERE id = $1", id,).Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.RoomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			g.logger.Warn("guest not found", "id", id)
			return nil, fmt.Errorf("guest with id %s not found", id)
		}
		g.logger.Error("error fetching guest by ID", "error", err)
        return nil, fmt.Errorf("error fetching guest by ID: %v", err)
	}
	g.logger.Info("successfully fetched guest", "id", id)
	return &guest, nil
}