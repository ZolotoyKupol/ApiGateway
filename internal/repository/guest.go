package repository

import (
	"apigateway/internal/models"
	"apigateway/internal/storage"
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"gorm.io/gorm"

)

type GuestRepo struct {
	store *storage.Storage
	logger *slog.Logger
}

func NewGuestRepo(store *storage.Storage, logger *slog.Logger) GuestRepoInterface {
	return &GuestRepo{store: store, logger: logger}
}

func (g *GuestRepo) GetGuestsRepo(ctx context.Context) ([]models.GuestDB, error) {
	var guests []models.GuestDB
	err := g.store.DB().WithContext(ctx).Find(&guests).Error
	if err != nil {
		g.logger.Error("failed to fetch guests", "error", err)
		return nil, errors.Wrap(err, "failed to fetch guests")
	}
	g.logger.Info("successfully fetched guests", "count", len(guests))
	return guests, nil
}

func (g *GuestRepo) CreateGuestRepo(ctx context.Context, guest models.GuestDB) (string, error) {
	err := g.store.DB().WithContext(ctx).Create(&guest).Error
	if err != nil {
		g.logger.Error("failed to create guest", "error", err)
		return "", errors.Wrap(err, "failed to create guest")
	}
	g.logger.Info("Successfully created guest", "id", guest.ID)
	return guest.ID, nil
}

func (g *GuestRepo) DeleteGuestRepo(ctx context.Context, id string) (error) {
	err := g.store.DB().WithContext(ctx).Where("id = ?", id).Delete(&models.GuestDB{}).Error
	if err != nil {
		g.logger.Error("failed to delete guest", "error", err)
		return errors.Wrap(err, "failed to delete guest")
	}
	g.logger.Info("successfully deleted guest", "id", id)
	return nil
}

func (g *GuestRepo) UpdateGuestRepo(ctx context.Context, id string, guest models.GuestDB) error {
	err := g.store.DB().WithContext(ctx).Model(&models.GuestDB{}).Where("id = ?", id).Updates(guest).Error
	if err != nil {
		g.logger.Error("failed to update guest", "error", err)
		return errors.Wrap(err, "failed to update guest")
	}
	g.logger.Info("guest updated successfully", "id", id)
	return nil
}

func (g *GuestRepo) GetGuestByIDRepo(ctx context.Context, id string) (*models.GuestDB, error) {
	var guest models.GuestDB
	err := g.store.DB().WithContext(ctx).First(&guest, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			g.logger.Warn("guest not found", "id", id)
			return nil, errors.Wrap(err, "guest not found")
		}
		g.logger.Error("error fetching guest by ID", "error", err)
        return nil, errors.Wrap(err, "error fetching guest by ID")
	}
	g.logger.Info("successfully fetched guest", "id", id)
	return &guest, nil
}

