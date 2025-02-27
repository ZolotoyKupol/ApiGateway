package repository

import (
	"apigateway/internal/models"
	"apigateway/internal/storage"
	"context"
	"github.com/pkg/errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type RoomRepo struct {
	store  *storage.Storage
	logger *slog.Logger
}

func NewRoomRepo(store *storage.Storage, logger *slog.Logger) RoomRepoInterface {
	return &RoomRepo{store: store, logger: logger}
}

func (r *RoomRepo) GetRoomsRepo(ctx context.Context) ([]models.RoomDB, error) {
	var rooms []models.RoomDB
	err := r.store.DB().WithContext(ctx).Find(&rooms).Error
	if err != nil {
		r.logger.Error("failed to fetch rooms", "error", err)
		return nil, errors.Wrap(err, "failed to fetch rooms")
	}
	r.logger.Info("successfully fetched rooms", "count", len(rooms))
	return rooms, nil
}

func (r *RoomRepo) CreateRoomRepo(ctx context.Context, room models.RoomDB) (string, error) {
	err := r.store.DB().WithContext(ctx).Create(&room).Error
	if err != nil {
		r.logger.Error("failed to create room", "error", err)
		return "", errors.Wrap(err, "failed to create room")
	}
	r.logger.Info("Successfully created room", "id", room.ID)
	return room.ID, nil
}

func (r *RoomRepo) UpdateRoomRepo(ctx context.Context, id string, room models.RoomDB) error {
	err := r.store.DB().WithContext(ctx).Model(&models.RoomDB{}).Where("id = ?", id).Updates(room).Error
	if err != nil {
		r.logger.Error("failed to update room", "error", err)
		return errors.Wrap(err, "failed to update room")
	}
	r.logger.Info("room updated successfully", "id", id)
	return nil
}

func (r *RoomRepo) DeleteRoomRepo(ctx context.Context, id string) error {
	err := r.store.DB().WithContext(ctx).Where("id = ?", id).Delete(&models.RoomDB{}).Error
	if err != nil {
		r.logger.Error("failed to delete room", "error", err)
		return errors.Wrap(err, "failed to delete room")
	}
	r.logger.Info("room deleted successfully", "id", id)
	return nil
}

func (r *RoomRepo) GetRoomByIDRepo(ctx context.Context, id string) (*models.RoomDB, error) {
	var room models.RoomDB
	err := r.store.DB().WithContext(ctx).Where("id = ?", id).First(&room).Error
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "room not found")
		}
		r.logger.Error("failed to fetch room", "error", err)
		return nil, errors.Wrap(err, "failed to fetch room")
	}
	r.logger.Info("room fetched successfully", "id", id)
	return &room, nil
}
