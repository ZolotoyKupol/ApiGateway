package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"github.com/pkg/errors"
	"log/slog"
)

type RoomUsecase struct {
	repo repository.RoomRepoInterface
	logger *slog.Logger
}

func NewRoomUsecase(repo repository.RoomRepoInterface, logger *slog.Logger) RoomUsecaseInteface{
	return &RoomUsecase{repo: repo, logger: logger}
}

func (u *RoomUsecase) GetRooms(ctx context.Context) ([]models.RoomResponse, error) {
	roomDB, err := u.repo.GetRoomsRepo(ctx)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching all rooms")
	}
	return models.ConvertToRoomResponseList(roomDB), nil
}

func (u *RoomUsecase) GetRoomByID(ctx context.Context, id string) (*models.RoomResponse, error) {
	roomDB, err := u.repo.GetRoomByIDRepo(ctx, id)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, errors.Wrap(err, "error fetching room")
	}
	roomResponse := roomDB.ConvertToRoomResponse()
	return &roomResponse, nil
}

func (u *RoomUsecase) CreateRoom(ctx context.Context, room models.RoomDB) (string, error) {
	return u.repo.CreateRoomRepo(ctx, room)
}

func (u *RoomUsecase) DeleteRoom(ctx context.Context, id string) error {
	return u.repo.DeleteRoomRepo(ctx, id)
}

func (u *RoomUsecase) UpdateRoom(ctx context.Context, id string, room models.RoomDB) error {
	return u.repo.UpdateRoomRepo(ctx, id, room)
}
