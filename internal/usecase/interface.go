package usecase

import (
	"apigateway/internal/models"
	"context"
)

type GuestUsecaseInterface interface {
	FetchAllGuests(ctx context.Context) ([]models.GuestResponse, error)
	CreateGuest(ctx context.Context, guest models.GuestDB) (string, error)
    DeleteGuest(ctx context.Context, id string) error
    UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error
    GetGuestByID(ctx context.Context, id string) (*models.GuestResponse, error)
}

type RoomUsecaseInteface interface {
	GetRooms(ctx context.Context) ([]models.RoomResponse, error)
	GetRoomByID(ctx context.Context, id string) (*models.RoomResponse, error)
	CreateRoom(ctx context.Context, room models.RoomDB) (string, error)
	DeleteRoom(ctx context.Context, id string) error
	UpdateRoom(ctx context.Context, id string, room models.RoomDB) error
}