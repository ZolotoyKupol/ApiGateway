package cache

import (
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"context"
	"sync"
)

type CachedRoom struct {
	roomRepo *repository.RoomRepo
	rooms    map[int]models.RoomDB
	mu 	 sync.RWMutex
}

func NewCachedRoom(roomRepo *repository.RoomRepo) *CachedRoom {
	return &CachedRoom{roomRepo: roomRepo, rooms: make(map[int]models.RoomDB)}
}


func (c *CachedRoom) GetAllRooms(ctx context.Context) ([]models.RoomDB, error) {
	c.mu.RLock()
	if len(c.rooms) > 0 {
		defer c.mu.RUnlock()
		var rooms []models.RoomDB
		for _, room := range c.rooms {
			rooms = append(rooms, room)
		}
		return rooms, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.rooms) == 0 {
		rooms, err := c.roomRepo.GetAllRooms(ctx)
		if err != nil {
			return nil, err
		}
		for _, room := range rooms {
			c.rooms[room.ID] = room
		}
	}
	
	var rooms []models.RoomDB
	for _, room := range c.rooms {
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (c *CachedRoom) GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error) {
	c.mu.RLock()
	if room, ok := c.rooms[id]; ok {
		defer c.mu.RUnlock()
		return &room, nil
	}
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()


	room, err := c.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.rooms[room.ID] = *room
	return room, nil
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


