package main

import (
	"apigateway/internal/app"
	"apigateway/internal/handlers"
	"apigateway/internal/repository"
	"apigateway/internal/storage"
	"apigateway/internal/usecase"
	"log/slog"
	"os"

)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		logger.Error("DATABASE_URL is not set")
		os.Exit(1)
	}
	store, err := storage.NewStorage(connString)
	if err != nil {
		logger.Error("initialization error db", "error", err)
		os.Exit(1)
	}
	defer func() {
		sqlDB, err := store.DB().DB()
		if err != nil {
			logger.Error("failed to get database connection", "error", err)
			return
		}
		closeErr := sqlDB.Close()
		if closeErr != nil {
			logger.Error("failed to close database connection", "error", closeErr)
		}
	}()

	guestRepo := repository.NewGuestRepo(store, logger)
	roomRepo := repository.NewRoomRepo(store, logger)

	guestUsecase := usecase.NewGuestUsecase(guestRepo, logger)
	roomUsecase := usecase.NewRoomUsecase(roomRepo, logger)

	guestHandler := handlers.NewHandlers(guestUsecase, logger)
	roomHandler := handlers.NewRoomHandler(roomUsecase, logger)


	router := app.SetupRouter(guestHandler, roomHandler)

	logger.Info("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to start server", "error", err)
        os.Exit(1)
	}
}