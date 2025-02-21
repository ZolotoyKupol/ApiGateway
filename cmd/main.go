package main

import (
	"apigateway/internal/handlers"
	"apigateway/internal/repository"
	"apigateway/internal/storage"
	"apigateway/internal/usecase"
	"context"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
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
		if closeErr := store.Conn().Close(context.Background()); closeErr != nil {
			logger.Error("failed to close database connection", "error", closeErr)
		}
	}()

	guestRepo := repository.NewGuestRepo(store, logger)
	guestUsecase := usecase.NewGuestUsecase(guestRepo, logger)
	handlers := handlers.NewHandlers(guestUsecase, logger)

	router := gin.Default()

	router.GET("/guests", handlers.FetchAllGuests)
	router.GET("/guests/:id", handlers.GetGuestByID)
	router.PUT("/guests/:id", handlers.UpdateGuest)
	router.POST("/guests", handlers.CreateGuest)
	router.DELETE("/guests/:id", handlers.DeleteGuest)

	logger.Info("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("failed to start server", "error", err)
        os.Exit(1)
	}
}