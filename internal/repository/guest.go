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

type GuestRepo struct {
	store  *storage.Storage
	logger *slog.Logger
}

func NewGuestRepo(store *storage.Storage, logger *slog.Logger) *GuestRepo {
	return &GuestRepo{store: store, logger: logger}
}

func (g *GuestRepo) GetAllGuests(ctx context.Context) ([]models.GuestDB, error) {
	var guests []models.GuestDB
	err := g.store.DB().WithContext(ctx).Find(&guests).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch guests")
	}
	if len(guests) == 0 {
		return nil, apperr.ErrNoData
	}
	g.logger.Debug("successfully fetched guests", "count", len(guests))
	return guests, nil
}

func (g *GuestRepo) CreateGuest(ctx context.Context, guest models.GuestDB) (int, error) {
	err := g.store.DB().WithContext(ctx).Create(&guest).Error
	if err != nil {
		g.logger.Error("failed to create guest", "error", err)
		return 0, errors.Wrap(err, "failed to create guest")
	}
	g.logger.Debug("Successfully created guest", "id", guest.ID)
	return guest.ID, nil
}

func (g *GuestRepo) DeleteGuest(ctx context.Context, id string) error {
	err := g.store.DB().WithContext(ctx).Where("id = ?", id).Delete(&models.GuestDB{}).Error
	if err != nil {
		g.logger.Error("failed to delete guest", "error", err)
		return errors.Wrap(err, "failed to delete guest")
	}
	g.logger.Debug("successfully deleted guest", "id", id)
	return nil
}

func (g *GuestRepo) UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error {
	err := g.store.DB().WithContext(ctx).Model(&models.GuestDB{}).Where("id = ?", id).Updates(guest).Error
	if err != nil {
		g.logger.Error("failed to update guest", "error", err)
		return errors.Wrap(err, "failed to update guest")
	}
	g.logger.Debug("guest updated successfully", "id", id)
	return nil
}

func (g *GuestRepo) GetGuestByID(ctx context.Context, id string) (*models.GuestDB, error) {
	var guest models.GuestDB
	err := g.store.DB().WithContext(ctx).First(&guest, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrNoData
		}
		g.logger.Error("error fetching guest by ID", "error", err)
		return nil, errors.Wrap(err, "error fetching guest by ID")
	}
	g.logger.Debug("successfully fetched guest", "id", id)
	return &guest, nil
}
