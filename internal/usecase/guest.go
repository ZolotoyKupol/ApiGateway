package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

var ErrNoData = errors.New("no data")

type GuestUC struct {
	repo   repository.GuestRepoProvider
	logger *slog.Logger
}

func NewGuestUsecase(repo repository.GuestRepoProvider, logger *slog.Logger) *GuestUC {
	return &GuestUC{repo: repo, logger: logger}
}

func (g *GuestUC) GetGuests(ctx context.Context) ([]models.GuestResponse, error) {
	guestDB, err := g.repo.GetAllGuests(ctx)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			return nil, ErrNoData
		}
		g.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching all guests")
	}
	return models.ConvertToGuestResponseList(guestDB), nil
}

func (g *GuestUC) CreateGuest(ctx context.Context, guest models.GuestDB) (int, error) {
	return g.repo.CreateGuest(ctx, guest)
}

func (g *GuestUC) DeleteGuest(ctx context.Context, id string) error {
	return g.repo.DeleteGuest(ctx, id)
}

func (g *GuestUC) UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error {
	return g.repo.UpdateGuest(ctx, id, guest)
}

func (g *GuestUC) GetGuestByID(ctx context.Context, id string) (*models.GuestResponse, error) {
	guestDB, err := g.repo.GetGuestByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			return nil, ErrNoData
		}
		g.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching guest")
	}
	guestResponse := guestDB.ConvertToGuestResponse()
	return &guestResponse, nil
}
