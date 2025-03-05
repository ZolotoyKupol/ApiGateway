package main

import (
	"apigateway/internal/app"
	"apigateway/internal/cache"
	"apigateway/internal/handlers"
	"apigateway/internal/repository"
	"apigateway/internal/storage"
	"apigateway/internal/usecase"
	"apigateway/migrations"
	"log/slog"
	"os"
)


func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)

	logger.Debug("logger initialized", "level", logLevel)

	logger.Debug("starting app", "env", os.Getenv("ENV"))
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		logger.Error("DATABASE_URL is not set")
		return
	}
	store, err := storage.NewStorage(connString)
	if err != nil {
		logger.Error("initialization error db", "error", err)
		return
	}
	logger.Debug("db connection string", "connString", connString)
	defer store.Close()
	
	guestRepo := repository.NewGuestRepo(store, logger)
	roomRepo := repository.NewRoomRepo(store, logger)
	
	roomCache := cache.NewCachedRoom(roomRepo)

	guestUsecase := usecase.NewGuestUsecase(guestRepo, logger)
	roomUsecase := usecase.NewRoomUsecase(roomCache, logger)

	guestHandler := handlers.NewHandlers(guestUsecase, logger)
	roomHandler := handlers.NewRoomHandler(roomUsecase, logger)


	router := app.SetupRouter(guestHandler, roomHandler)

	if err := migrations.Migrate(connString); err != nil {
		logger.Error("migration error", "error", err)
		return
	}
	logger.Info("migrations applied")

	logger.Info("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to start server", "error", err)
        os.Exit(1)
	}
}