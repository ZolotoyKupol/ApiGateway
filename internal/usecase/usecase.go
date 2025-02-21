package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"errors"
	"log/slog"

	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrNoData = errors.New("no data")

type GuestUsecase struct {
	repo *repository.GuestRepo
	logger *slog.Logger
}

func NewGuestUsecase(repo *repository.GuestRepo, logger *slog.Logger) *GuestUsecase {
	return  &GuestUsecase{repo: repo, logger: logger}
}


func (g *GuestUsecase) FetchAllGuests(ctx context.Context) ([]models.Guest, error) {
	guests, err := g.repo.GetGuestsRepo(ctx) 
		if err != nil {
			g.logger.Error("failed to fetch guests", "error", err)
			return nil, fmt.Errorf("failed to fetch guests: %v", err)
		}
		g.logger.Info("successfully fetched all guests", "count", len(guests))
		return guests, nil
	}
	

func (g *GuestUsecase) CreateGuest(ctx context.Context, guest models.Guest) (string, error) {
	id, err := g.repo.CreateGuestRepo(ctx, guest)
	if err != nil {
		g.logger.Error("failed to create guest", "error", err)
        return "", fmt.Errorf("failed to create guest: %v", err)
    }
    g.logger.Info("guest created successfully", "id", id)
    return id, nil
}

func (g *GuestUsecase) DeleteGuest(ctx context.Context, id string) error {
	err := g.repo.DeleteGuestRepo(ctx, id)
	if err != nil {
		g.logger.Error("error deleting guest", "error", err)
		return fmt.Errorf("error deleting guest: %v", err)
	}
	g.logger.Info("guest deleted successfully", "id", id)
	return nil
}

func (g *GuestUsecase) UpdateGuest(ctx context.Context, id string, guest models.Guest) error {
	err := g.repo.UpdateGuestRepo(ctx, id, guest)
	if err != nil {
		g.logger.Error("error updating guest", "error", err)
		return fmt.Errorf("error updating guest: %v", err)
	}
	g.logger.Info("guest updated successfully", "id", id)
	return nil
}

func (g *GuestUsecase) GetGuestByID(ctx context.Context, id string) (*models.Guest, error) {
	guest, err := g.repo.GetGuestByIDRepo(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoData
		}
		return nil, fmt.Errorf("failed to fetch guest by ID: %v", err)
    }
    g.logger.Info("successfully fetched guest", "id", id)
    return guest, nil
}
