package app

import (
	"apigateway/internal/handlers"
	"apigateway/internal/metrics"
	"apigateway/internal/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(guestHandler handlers.GuestProvider, roomHandler handlers.RoomProvider, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(metrics.PrometheusMiddleware(), middleware.LoggingMiddleware(logger))

	guestGroup := router.Group("/guests")
	guestGroup.GET("", guestHandler.GetAllGuests)
	guestGroup.GET("/:id", guestHandler.GetGuestByID)
	guestGroup.POST("", guestHandler.CreateGuest)
	guestGroup.PUT("/:id", guestHandler.UpdateGuest)
	guestGroup.DELETE("/:id", guestHandler.DeleteGuest)

	roomGroup := router.Group("/rooms")
	roomGroup.GET("", roomHandler.GetAllRooms)
	roomGroup.GET("/:id", roomHandler.GetRoomByID)
	roomGroup.POST("", roomHandler.CreateRoom)
	roomGroup.PUT("/:id", roomHandler.UpdateRoom)
	roomGroup.DELETE("/:id", roomHandler.DeleteRoom)

	return router
}
