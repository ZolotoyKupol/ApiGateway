package app

import (
	"apigateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(guestHandler handlers.GuestProvider, roomHandler handlers.RoomProvider) *gin.Engine {
	router := gin.Default()

	guestGroup := router.Group("/guests")
    {
        guestGroup.GET("", guestHandler.FetchAllGuests)
        guestGroup.GET("/:id", guestHandler.GetGuestByID)
        guestGroup.POST("", guestHandler.CreateGuest)
        guestGroup.PUT("/:id", guestHandler.UpdateGuest)
        guestGroup.DELETE("/:id", guestHandler.DeleteGuest)
    }

	roomGroup := router.Group("/rooms")
    {
        roomGroup.GET("", roomHandler.GetRoomsHandler)
        roomGroup.GET("/:id", roomHandler.GetRoomByIDHandler)
        roomGroup.POST("", roomHandler.CreateRoomHandler)
        roomGroup.PUT("/:id", roomHandler.UpdateRoomHandler)
        roomGroup.DELETE("/:id", roomHandler.DeleteRoomHandler)
    }
	return router
}