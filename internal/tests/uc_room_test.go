package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"apigateway/internal/apperr"
	"apigateway/internal/models"
	"apigateway/internal/repository"
	"apigateway/internal/storage"
	"apigateway/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRoomUsecase(t *testing.T) {
	connString := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.RoomDB{})
	require.NoError(t, err)
	defer func() {
		if err := db.Migrator().DropTable(&models.RoomDB{}); err != nil {
			t.Fatalf("failed to drop table: %v", err)
		}
	}()

	storage, err := storage.NewStorage(connString)
	require.NoError(t, err)
	defer storage.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	repo := repository.NewRoomRepo(storage, logger)

	uc := usecase.NewRoomUsecase(repo, logger)

	ctx := context.Background()

	room := models.RoomDB{
		Number:     "101",
		Floor:      "1",
		RoomSize:   20.0,
		Status:     "available",
		OccupiedBy: "",
		CheckIn:    time.Time{},
		CheckOut:   time.Time{},
	}

	t.Run("CreateRoom", func(t *testing.T) {
		id, err := uc.CreateRoom(ctx, room)
		require.NoError(t, err)
		assert.NotZero(t, id)

		createdRoom, err := uc.GetRoomByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, room.Number, createdRoom.Number)
		assert.Equal(t, room.Floor, createdRoom.Floor)
		assert.InEpsilon(t, room.RoomSize, createdRoom.RoomSize, 0.0001)
		assert.Equal(t, room.Status, createdRoom.Status)
		assert.Equal(t, room.OccupiedBy, createdRoom.OccupiedBy)
	})

	t.Run("GetAllRooms", func(t *testing.T) {
		rooms, err := uc.GetRooms(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, rooms)
	})

	t.Run("GetRoomByID", func(t *testing.T) {
		id, err := uc.CreateRoom(ctx, room)
		require.NoError(t, err)

		room, err := uc.GetRoomByID(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, room)
		assert.Equal(t, "101", room.Number)
		assert.Equal(t, "1", room.Floor)
		assert.InEpsilon(t, 20.0, room.RoomSize, 0.0001)
		assert.Equal(t, "available", room.Status)
		assert.Equal(t, "", room.OccupiedBy)
	})

	t.Run("UpdateRoom", func(t *testing.T) {
		id, err := uc.CreateRoom(ctx, room)
		require.NoError(t, err)

		updatedData := models.RoomDB{
			Number:     "102",
			Floor:      "2",
			RoomSize:   30.0,
			Status:     "occupied",
			OccupiedBy: "John Doe",
			CheckIn:    time.Now(),
			CheckOut:   time.Now().Add(24 * time.Hour),
		}

		err = uc.UpdateRoom(ctx, id, updatedData)
		require.NoError(t, err)

		updatedRoom, err := uc.GetRoomByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, updatedData.Number, updatedRoom.Number)
		assert.Equal(t, updatedData.Floor, updatedRoom.Floor)
		assert.InEpsilon(t, updatedData.RoomSize, updatedRoom.RoomSize, 0.0001)
		assert.Equal(t, updatedData.Status, updatedRoom.Status)
		assert.Equal(t, updatedData.OccupiedBy, updatedRoom.OccupiedBy)
	})

	t.Run("DeleteRoom", func(t *testing.T) {
		id, err := uc.CreateRoom(ctx, room)
		require.NoError(t, err)

		err = uc.DeleteRoom(ctx, id)
		require.NoError(t, err)
		_, err = uc.GetRoomByID(ctx, id)
		assert.ErrorIs(t, err, apperr.ErrNoData)
	})
}
