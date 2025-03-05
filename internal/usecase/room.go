package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type RoomUC struct {
	repo   repository.RoomRepoProvider
	logger *slog.Logger
}

func NewRoomUsecase(repo repository.RoomRepoProvider, logger *slog.Logger) *RoomUC {
	return &RoomUC{repo: repo, logger: logger}
}

func (u *RoomUC) GetRooms(ctx context.Context) ([]models.RoomResponse, error) {
	roomDB, err := u.repo.GetAllRooms(ctx)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			return nil, ErrNoData
		}
		u.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching all rooms")
	}
	return models.ConvertToRoomResponseList(roomDB), nil
}

func (u *RoomUC) GetRoomByID(ctx context.Context, id int) (*models.RoomResponse, error) {
	roomDB, err := u.repo.GetRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNoData) {
			return nil, ErrNoData
		}
		u.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching room")
	}
	roomResponse := roomDB.ConvertToRoomResponse()
	return &roomResponse, nil
}

func (u *RoomUC) CreateRoom(ctx context.Context, room models.RoomDB) (int, error) {
	return u.repo.CreateRoom(ctx, room)
}

func (u *RoomUC) DeleteRoom(ctx context.Context, id int) error {
	return u.repo.DeleteRoom(ctx, id)
}

func (u *RoomUC) UpdateRoom(ctx context.Context, id int, room models.RoomDB) error {
	return u.repo.UpdateRoom(ctx, id, room)
}
