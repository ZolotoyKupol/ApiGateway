package repository

import (
	"apigateway/internal/models"
	"apigateway/internal/storage"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)


type RoomRepo struct {
	store *storage.Storage
	logger *slog.Logger
}

func NewRoomRepo(store *storage.Storage, logger *slog.Logger) *RoomRepo{
	return &RoomRepo{store: store, logger: logger}
}

func (r *RoomRepo) GetRoomsRepo(ctx context.Context) ([]models.Room, error) {
	rows, err := r.store.Query("SELECT id, number, floor, room_size, status, occupied_by, check_in, check_out FROM rooms")
	if err != nil {
		r.logger.Error("failed to fetch all rooms", "error", err)
        return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Number, &room.Floor, &room.RoomSize, &room.Status, &room.OccupiedBy, &room.CheckIn, &room.CheckOut)
		if err != nil {
			r.logger.Error("failed to scan room", "error", err)
			return nil, err
		}
		rooms = append(rooms, room)
	}
	r.logger.Info("all rooms retrieved successfully", "count", len(rooms))
	return rooms, nil
}

func (r *RoomRepo) CreateRoomRepo(ctx context.Context, room models.Room) (string, error) {
	query := `INSERT INTO rooms (id, number, floor, room_size, status, occupied_by, check_in, check_out) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id string
	err := r.store.Conn().QueryRow(ctx, query, room.ID, room.Number, room.Floor, room.RoomSize, room.Status, room.OccupiedBy, room.CheckIn, room.CheckOut,).Scan(&id)
	if err != nil {
		r.logger.Error("error creating room", "error", err)
		return "", err
	}
	return id, nil
}

func (r *RoomRepo) UpdateRoomRepo(ctx context.Context, id string, room models.Room) error {
	query := `UPDATE rooms SET number = $1, floor = $2, room_size = $3, status = $4, occupied_by = $5, check_in = $6, check_out = $7 WHERE id = $8`
	_, err := r.store.Conn().Exec(ctx, query, room.Number, room.Floor, room.RoomSize, room.Status, room.OccupiedBy, room.CheckIn, room.CheckOut, room.ID)
	if err != nil {
		r.logger.Error("failed to update room", "error", err)
		return err
	}
	r.logger.Info("room update successfully", "room_id", room.ID)
	return nil
}


func (r *RoomRepo) DeleteRoomRepo(ctx context.Context, id string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.store.Conn().Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete room", "error", err)
		return err
	}
	r.logger.Info("room delete successfully", "id", id)
	return nil
}

func (r *RoomRepo) GetRoomByIDRepo(ctx context.Context, id string) (*models.Room, error) {
	query := `SELECT id, number, floor, room_size, status, occupied_by, check_in, check_out FROM rooms WHERE id = $1`
	var room models.Room
	err := r.store.Conn().QueryRow(ctx, query, id).Scan(&room.ID, &room.Number, &room.Floor, &room.RoomSize, &room.Status, &room.OccupiedBy, &room.CheckIn, &room.CheckOut,)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("room not found", "id", id)
			return nil, nil
		}
		r.logger.Error("failed to get room by ID", "error", err)
		return nil, err
	}
	r.logger.Info("successfully fetched room", "id", id)
	return &room, nil
}