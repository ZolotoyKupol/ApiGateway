package main

import (
	"apigateway/internal/app"
	"apigateway/internal/cache"
	"apigateway/internal/handlers"
	"apigateway/internal/metrics"
	"apigateway/internal/repository"
	"apigateway/internal/storage"
	"apigateway/internal/tracing"
	"apigateway/internal/usecase"
	"apigateway/migrations"
	"fmt"
	"log/slog"
	"os"

	"github.com/pkg/errors"
)

func loggerInit(level string) (*slog.Logger, error) {
	logLevel := slog.LevelInfo
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, errors.Wrap(err, "failed to parse log level")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.Info(fmt.Sprintf("minimum log level set: %s", logLevel.String()))
	return logger, nil
}

func main() {
	logger, err := loggerInit(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logger.Error("failed to initialize logger", "error", err)
		return
	}

	logger.Debug("logger initialized")

	shutdown := tracing.InitTracer("api-gateway")
	defer shutdown()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		logger.Error("DATABASE_URL is not set")
		return
	}

	if err := migrations.Migrate(connString); err != nil {
		logger.Error("migration error", "error", err)
		return
	}
	logger.Info("migrations applied")

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

	metrics.Init(os.Getenv("METRICS_PORT"))

	router := app.SetupRoutes(guestHandler, roomHandler, logger)

	logger.Info("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to start server", "error", err)
		return
	}
}
