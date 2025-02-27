package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"github.com/pkg/errors"
	"log/slog"
)

var ErrNoData = errors.New("no data")

type GuestUsecase struct {
	repo repository.GuestRepoInterface
	logger *slog.Logger
}

func NewGuestUsecase(repo repository.GuestRepoInterface, logger *slog.Logger) GuestUsecaseInterface {
	return  &GuestUsecase{repo: repo, logger: logger}
}


func (g *GuestUsecase) FetchAllGuests(ctx context.Context) ([]models.GuestResponse, error) {
	guestDB, err := g.repo.GetGuestsRepo(ctx)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching all guests")
	}
	return models.ConvertToGuestResponseList(guestDB), nil
}
	

func (g *GuestUsecase) CreateGuest(ctx context.Context, guest models.GuestDB) (string, error) {
	return g.repo.CreateGuestRepo(ctx, guest)
}

func (g *GuestUsecase) DeleteGuest(ctx context.Context, id string) error {
	return g.repo.DeleteGuestRepo(ctx, id)
}

func (g *GuestUsecase) UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error {
	return g.repo.UpdateGuestRepo(ctx, id, guest)
}

func (g *GuestUsecase) GetGuestByID(ctx context.Context, id string) (*models.GuestResponse, error) {
	guestDB, err := g.repo.GetGuestByIDRepo(ctx, id)
	if err != nil {
		g.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching guest")
	}
	guestResponse := guestDB.ConvertToGuestResponse()
	return &guestResponse, nil
}

