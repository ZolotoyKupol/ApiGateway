package cache

import (
	"apigateway/internal/metrics"
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"log/slog"
	"sync"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

type CachedRoom struct {
	roomRepo repository.RoomProvider
	mu       sync.RWMutex
	rooms    map[int]models.RoomDB
}

func NewCachedRoom(roomRepo repository.RoomProvider) *CachedRoom {
	return &CachedRoom{roomRepo: roomRepo, rooms: make(map[int]models.RoomDB)}
}

func (c *CachedRoom) GetAllRooms(ctx context.Context) ([]models.RoomDB, error) {
	ctx, span := otel.Tracer("CachedRoom").Start(ctx, "CachedRoom.GetAllRooms")
	defer span.End()

	return c.roomRepo.GetAllRooms(ctx)
}

func (c *CachedRoom) GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error) {
	ctx, span := otel.Tracer("CachedRoom").Start(ctx, "CachedRoom.GetRoomByID")
	defer span.End()

	if room, ok := c.Get(ctx, id); ok {
		return room, nil
	}

	dbRoom, err := c.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "not found room in cache")
	}

	if err := c.Set(ctx, *dbRoom); err != nil {
		slog.Debug("failed to set room in cache", "err", err)
	}
	return dbRoom, nil
}

func (c *CachedRoom) CreateRoom(ctx context.Context, room models.RoomDB) (int, error) {
	ctx, span := otel.Tracer("CachedRoom").Start(ctx, "CachedRoom.CreateRoom")
	defer span.End()

	id, err := c.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return 0, err
	}

	room.ID = id

	if err := c.Set(ctx, room); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *CachedRoom) DeleteRoom(ctx context.Context, id int) error {
	ctx, span := otel.Tracer("CachedRoom").Start(ctx, "CachedRoom.DeleteRoom")
	defer span.End()

	err := c.roomRepo.DeleteRoom(ctx, id)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.rooms, id)
	metrics.UpdateCacheSizeMetric(c.rooms)
	return nil
}

func (c *CachedRoom) UpdateRoom(ctx context.Context, id int, room models.RoomDB) error {
	ctx, span := otel.Tracer("CachedRoom").Start(ctx, "CachedRoom.UpdateRoom")
	defer span.End()

	err := c.roomRepo.UpdateRoom(ctx, id, room)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[id] = room
	metrics.UpdateCacheSizeMetric(c.rooms)
	return nil
}

func (c *CachedRoom) Get(ctx context.Context, id int) (*models.RoomDB, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if room, ok := c.rooms[id]; ok {
		return &room, true
	}
	return nil, false
}

func (c *CachedRoom) Set(ctx context.Context, room models.RoomDB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[room.ID] = room
	metrics.UpdateCacheSizeMetric(c.rooms)
	return nil
}

func (c *CachedRoom) SetAll(ctx context.Context, rooms []models.RoomDB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, room := range rooms {
		c.rooms[room.ID] = room
	}
	return nil
}
