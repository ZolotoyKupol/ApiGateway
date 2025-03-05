package app

import (
	"apigateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(guestHandler handlers.GuestProvider, roomHandler handlers.RoomProvider) *gin.Engine {
	router := gin.Default()

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
