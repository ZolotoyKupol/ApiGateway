package usecase

import (
	"apigateway/internal/models"
	"context"
)

type GuestProvider interface {
	GetGuests(ctx context.Context) ([]models.GuestDB, error)
	CreateGuest(ctx context.Context, guest models.GuestDB) (int, error)
	DeleteGuest(ctx context.Context, id string) error
	UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error
	GetGuestByID(ctx context.Context, id string) (*models.GuestDB, error)
}

type RoomProvider interface {
	GetRooms(ctx context.Context) ([]models.RoomDB, error)
	GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error)
	CreateRoom(ctx context.Context, room models.RoomDB) (int, error)
	DeleteRoom(ctx context.Context, id int) error
	UpdateRoom(ctx context.Context, id int, room models.RoomDB) error
}
