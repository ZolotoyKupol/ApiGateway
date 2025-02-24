package usecase

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"log/slog"

)

type RoomUsecase struct {
	repo *repository.RoomRepo
	logger *slog.Logger
}

func NewRoomUsecase(repo *repository.RoomRepo, logger *slog.Logger) *RoomUsecase{
	return &RoomUsecase{repo: repo, logger: logger}
}

func (u *RoomUsecase) GetRooms(ctx context.Context) ([]models.Room, error) {
	return u.repo.GetRoomsRepo(ctx)
}

func (u *RoomUsecase) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	return u.repo.GetRoomByIDRepo(ctx, id)
}

func (u *RoomUsecase) CreateRoom(ctx context.Context, room models.Room) (string, error) {
	return u.repo.CreateRoomRepo(ctx, room)
}

func (u *RoomUsecase) DeleteRoom(ctx context.Context, id string) error {
	return u.repo.DeleteRoomRepo(ctx, id)
}

func (u *RoomUsecase) UpdateRoom(ctx context.Context, id string, room models.Room) error {
	return u.repo.UpdateRoomRepo(ctx, id, room)
}
