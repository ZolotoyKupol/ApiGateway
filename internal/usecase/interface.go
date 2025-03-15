package usecase

import (
	"apigateway/internal/models"
	"context"
)

type GuestUCProvider interface {
	GetGuests(ctx context.Context) ([]models.GuestResponse, error)
	CreateGuest(ctx context.Context, guest models.GuestDB) (int, error)
	DeleteGuest(ctx context.Context, id string) error
	UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error
	GetGuestByID(ctx context.Context, id string) (*models.GuestResponse, error)
}

type RoomUCProvider interface {
	GetRooms(ctx context.Context) ([]models.RoomResponse, error)
	GetRoomByID(ctx context.Context, id int) (*models.RoomResponse, error)
	CreateRoom(ctx context.Context, room models.RoomDB) (int, error)
	DeleteRoom(ctx context.Context, id int) error
	UpdateRoom(ctx context.Context, id int, room models.RoomDB) error
}
