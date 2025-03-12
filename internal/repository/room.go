package repository

import (
	"apigateway/internal/apperr"
	"apigateway/internal/models"
	"apigateway/internal/storage"
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/gorm"

)

type RoomRepo struct {
	store  *storage.Storage
	logger *slog.Logger
}

func NewRoomRepo(store *storage.Storage, logger *slog.Logger) *RoomRepo {
	return &RoomRepo{store: store, logger: logger}
}

func (r *RoomRepo) GetAllRooms(ctx context.Context) ([]models.RoomDB, error) {
	var rooms []models.RoomDB
	err := r.store.DB().WithContext(ctx).Find(&rooms).Error
	if err != nil {
		r.logger.Error("failed to fetch rooms", "error", err)
		return nil, errors.Wrap(err, "failed to fetch rooms")
	}
	r.logger.Debug("successfully fetched rooms", "count", len(rooms))
	return rooms, nil
}

func (r *RoomRepo) CreateRoom(ctx context.Context, room models.RoomDB) (int, error) {
	err := r.store.DB().WithContext(ctx).Create(&room).Error
	if err != nil {
		r.logger.Error("failed to create room", "error", err)
		return 0, errors.Wrap(err, "failed to create room")
	}
	r.logger.Debug("Successfully created room", "id", room.ID)
	return room.ID, nil
}

func (r *RoomRepo) UpdateRoom(ctx context.Context, id int, room models.RoomDB) error {
	err := r.store.DB().WithContext(ctx).Model(&models.RoomDB{}).Where("id = ?", id).Updates(room).Error
	if err != nil {
		r.logger.Error("failed to update room", "error", err)
		return errors.Wrap(err, "failed to update room")
	}
	r.logger.Debug("room updated successfully", "id", id)
	return nil
}

func (r *RoomRepo) DeleteRoom(ctx context.Context, id int) error {
	err := r.store.DB().WithContext(ctx).Where("id = ?", id).Delete(&models.RoomDB{}).Error
	if err != nil {
		r.logger.Error("failed to delete room", "error", err)
		return errors.Wrap(err, "failed to delete room")
	}
	r.logger.Debug("room deleted successfully", "id", id)
	return nil
}

func (r *RoomRepo) GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error) {
	var room models.RoomDB
	err := r.store.DB().WithContext(ctx).Where("id = ?", id).First(&room).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("room not found", "id", id)
			return nil, apperr.ErrNoData
		}
		r.logger.Error("failed to fetch room", "error", err)
		return nil, errors.Wrap(err, "failed to fetch room")
	}
	r.logger.Debug("room fetched successfully", "id", id)
	return &room, nil
}
