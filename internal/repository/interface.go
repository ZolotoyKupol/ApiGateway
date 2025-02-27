package repository

import (
	"apigateway/internal/models"
	"context"
)


type GuestRepoInterface interface {
	GetGuestsRepo(ctx context.Context) ([]models.GuestDB, error)
	CreateGuestRepo(ctx context.Context, guest models.GuestDB) (string, error)
	DeleteGuestRepo(ctx context.Context, id string) (error)
	UpdateGuestRepo(ctx context.Context, id string, guest models.GuestDB) error
	GetGuestByIDRepo(ctx context.Context, id string) (*models.GuestDB, error)
}


type RoomRepoInterface interface {
	CreateRoomRepo(ctx context.Context, room models.RoomDB) (string, error)
	DeleteRoomRepo(ctx context.Context, id string) error
	GetRoomByIDRepo(ctx context.Context, id string) (*models.RoomDB, error)
	GetRoomsRepo(ctx context.Context) ([]models.RoomDB, error)
	UpdateRoomRepo(ctx context.Context, id string, room models.RoomDB) error
}