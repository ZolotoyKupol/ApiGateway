package cache

import (
	"apigateway/internal/mocks"
	"apigateway/internal/models"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubRoomProvider struct {
	rooms []models.RoomDB
}

func (s *StubRoomProvider) GetAllRooms(ctx context.Context) ([]models.RoomDB, error) {
	return s.rooms, nil
}

func (s *StubRoomProvider) GetRoomByID(ctx context.Context, id int) (*models.RoomDB, error) {
	if id == 1 {
		return &models.RoomDB{ID: 1, Number: "101"}, nil
	}
	return nil, fmt.Errorf("room not found")
}

func (s *StubRoomProvider) CreateRoom(ctx context.Context, room models.RoomDB) (int, error) {
	return 1, nil
}

func (s *StubRoomProvider) DeleteRoom(ctx context.Context, id int) error {
	if id == 1 {
		return nil
	}
	return fmt.Errorf("room not found")
}

func (s *StubRoomProvider) UpdateRoom(ctx context.Context, id int, room models.RoomDB) error {
	if id == 1 {
		return nil
	}
	return fmt.Errorf("room not found")
}

func TestGetAllRooms(t *testing.T) {
	testCases := []struct {
		name        string
		stubRepo    *StubRoomProvider
		expectedLen int
	}{
		{
			name: "валидный тест",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{
					{ID: 1, Number: "101"},
					{ID: 2, Number: "102"},
				},
			},
			expectedLen: 2,
		},
		{
			name: "пустой кэш",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCachedRoom(tc.stubRepo)
			rooms, err := cache.GetAllRooms(context.Background())
			require.NoError(t, err)
			assert.Len(t, rooms, tc.expectedLen)
		})
	}
}

func TestGetRoomByID(t *testing.T) {
	testCases := []struct {
		name        string
		stubRepo    *StubRoomProvider
		id          int
		expectedErr bool
	}{
		{
			name: "комната найдена в кэше",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{
					{ID: 1, Number: "101"},
				},
			},
			id:          1,
			expectedErr: false,
		},
		{
			name: "комната найдена в репозитории",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			id:          1,
			expectedErr: false,
		},
		{
			name: "комната не найдена",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			id:          999,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCachedRoom(tc.stubRepo)
			room, err := cache.GetRoomByID(context.Background(), tc.id)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.id, room.ID)
			}
		})
	}
}

func TestCreateRoom(t *testing.T) {
	testCases := []struct {
		name        string
		stubRepo    *StubRoomProvider
		room        models.RoomDB
		expectedID  int
		expectedErr bool
	}{
		{
			name: "успешное создание комнаты",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			room:        models.RoomDB{Number: "201"},
			expectedID:  1,
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCachedRoom(tc.stubRepo)
			id, err := cache.CreateRoom(context.Background(), tc.room)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)

				cache.mu.RLock()
				defer cache.mu.RUnlock()

				roomInCache, ok := cache.rooms[id]
				require.True(t, ok)
				assert.Equal(t, tc.room.Number, roomInCache.Number)
			}
		})
	}
}

func TestDeleteRoom(t *testing.T) {
	testCases := []struct {
		name        string
		stubRepo    *StubRoomProvider
		id          int
		expectedErr bool
	}{
		{
			name: "успешное удаление",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{
					{ID: 1, Number: "101"},
				},
			},
			id:          1,
			expectedErr: false,
		},
		{
			name: "комната не найдена",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			id:          999,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCachedRoom(tc.stubRepo)
			err := cache.DeleteRoom(context.Background(), tc.id)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				cache.mu.RLock()
				defer cache.mu.RUnlock()

				_, ok := cache.rooms[tc.id]
				assert.False(t, ok)
			}
		})
	}
}

func TestUpdateRoom(t *testing.T) {
	testCases := []struct {
		name        string
		stubRepo    *StubRoomProvider
		id          int
		newRoom     models.RoomDB
		expectedErr bool
	}{
		{
			name: "успешное обновление",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{
					{ID: 1, Number: "101"},
				},
			},
			id:          1,
			newRoom:     models.RoomDB{ID: 1, Number: "UpdatedRoom"},
			expectedErr: false,
		},
		{
			name: "комната не найдена",
			stubRepo: &StubRoomProvider{
				rooms: []models.RoomDB{},
			},
			id:          999,
			newRoom:     models.RoomDB{},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := NewCachedRoom(tc.stubRepo)
			err := cache.UpdateRoom(context.Background(), tc.id, tc.newRoom)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				cache.mu.RLock()
				defer cache.mu.RUnlock()

				roomInCache, ok := cache.rooms[tc.id]
				require.True(t, ok)
				assert.Equal(t, tc.newRoom.Number, roomInCache.Number)
			}
		})
	}
}

func TestDeleteRoom_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomProvider(ctrl)

	mockRepo.EXPECT().
		DeleteRoom(gomock.Any(), 1).
		Return(nil)

	cache := NewCachedRoom(mockRepo)
	err := cache.DeleteRoom(context.Background(), 1)

	require.NoError(t, err)
}

func TestCreateRoom_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomProvider(ctrl)

	mockRepo.EXPECT().
		CreateRoom(gomock.Any(), models.RoomDB{Number: "201"}).
		Return(1, nil)

	cache := NewCachedRoom(mockRepo)
	id, err := cache.CreateRoom(context.Background(), models.RoomDB{Number: "201"})

	require.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestGetRoomByID_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomProvider(ctrl)

	mockRepo.EXPECT().
		GetRoomByID(gomock.Any(), 1).
		Return(&models.RoomDB{ID: 1, Number: "101"}, nil)

	cache := NewCachedRoom(mockRepo)

	room, err := cache.GetRoomByID(context.Background(), 1)

	require.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, "101", room.Number)
}

func TestUpdateRoom_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomProvider(ctrl)

	mockRepo.EXPECT().
		UpdateRoom(gomock.Any(), 1, models.RoomDB{ID: 1, Number: "UpdatedRoom"}).
		Return(nil)

	cache := NewCachedRoom(mockRepo)

	if err := cache.Set(context.Background(), models.RoomDB{ID: 1, Number: "101"}); err != nil {
		t.Fatal("failed to set room in cache", err)
	}

	err := cache.UpdateRoom(context.Background(), 1, models.RoomDB{ID: 1, Number: "UpdatedRoom"})

	require.NoError(t, err)

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	roomInCache, ok := cache.rooms[1]
	require.True(t, ok)
	assert.Equal(t, "UpdatedRoom", roomInCache.Number)
}

func TestGetAllRooms_WithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRoomProvider(ctrl)

	mockRepo.EXPECT().
		GetAllRooms(gomock.Any()).
		Return([]models.RoomDB{
			{ID: 1, Number: "101"},
			{ID: 2, Number: "102"},
		}, nil)

	cache := NewCachedRoom(mockRepo)

	rooms, err := cache.GetAllRooms(context.Background())

	require.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, "101", rooms[0].Number)
	assert.Equal(t, "102", rooms[1].Number)

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	assert.Len(t, cache.rooms, 2)
	assert.Equal(t, "101", cache.rooms[1].Number)
	assert.Equal(t, "102", cache.rooms[2].Number)
}
