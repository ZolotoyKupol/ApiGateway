package cache

import (
	"apigateway/internal/apperr"
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"apigateway/internal/metrics"
	"context"
	"sync"
)

type CachedRoom struct {
	roomRepo repository.RoomProvider
	mu 	 sync.RWMutex
	rooms    map[int]models.RoomDB
	
}

func NewCachedRoom(roomRepo repository.RoomProvider) *CachedRoom {
	return &CachedRoom{roomRepo: roomRepo, rooms: make(map[int]models.RoomDB)}
}


func (c *CachedRoom) GetAllRooms(ctx context.Context) ([]models.RoomDB, error) {
	c.rooms = make(map[int]models.RoomDB)
	if len(c.rooms) == 0 {
		rooms, err := c.roomRepo.GetAllRooms(ctx)
		if err != nil {
			return nil, err
		}
		if err := c.SetAll(ctx, rooms); err != nil {
			return nil, err
		}
	}
	metrics.UpdateCacheSizeMetric(len(c.rooms))
	return c.convertMapToSlice(), nil
}

func (c *CachedRoom) GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error) {
	room, err := c.Get(ctx, id)
	if err == nil {
		return room, nil
	}

	dbRoom, err := c.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.Set(ctx, *dbRoom)
	return dbRoom, nil
}

func (c *CachedRoom) CreateRoom(ctx context.Context, room models.RoomDB) (int, error) {
	id, err := c.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		return 0, err
	}

	room.ID = id

	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[id] = room
	return id, nil
}

func (c *CachedRoom) DeleteRoom(ctx context.Context, id int) error {
	err := c.roomRepo.DeleteRoom(ctx, id)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.rooms, id)
	return nil
}

func (c *CachedRoom) UpdateRoom(ctx context.Context, id int, room models.RoomDB) error {
	err := c.roomRepo.UpdateRoom(ctx, id, room)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[id] = room
	return nil
}

func (c *CachedRoom) Get(ctx context.Context, id int) (*models.RoomDB, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if room, ok := c.rooms[id]; ok {
		return &room, nil
	}
	return nil, apperr.ErrNoDataCache
}

func (c *CachedRoom) Set(ctx context.Context, room models.RoomDB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rooms[room.ID] = room
	return nil
}

func (c *CachedRoom) convertMapToSlice() []models.RoomDB {
	rooms := make([]models.RoomDB, 0, len(c.rooms))
	for _, room := range c.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (c *CachedRoom) SetAll(ctx context.Context, rooms []models.RoomDB) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, room := range rooms {
		c.rooms[room.ID] = room
	}
	return nil
}
