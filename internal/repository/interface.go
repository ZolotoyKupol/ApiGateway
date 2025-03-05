package repository

import (
	"apigateway/internal/models"
	"context"
)

type GuestRepoProvider interface {
	GetAllGuests(ctx context.Context) ([]models.GuestDB, error)
	CreateGuest(ctx context.Context, guest models.GuestDB) (int, error)
	DeleteGuest(ctx context.Context, id string) error
	UpdateGuest(ctx context.Context, id string, guest models.GuestDB) error
	GetGuestByID(ctx context.Context, id string) (*models.GuestDB, error)
}

type RoomRepoProvider interface {
	CreateRoom(ctx context.Context, room models.RoomDB) (int, error)
	DeleteRoom(ctx context.Context, id int) error
	GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error)
	GetAllRooms(ctx context.Context) ([]models.RoomDB, error)
	UpdateRoom(ctx context.Context, id int, room models.RoomDB) error
}
