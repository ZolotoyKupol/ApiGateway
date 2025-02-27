package handlers

import "github.com/gin-gonic/gin"

type RoomProvider interface {
	GetRoomsHandler(c *gin.Context)
	CreateRoomHandler(c *gin.Context)
	UpdateRoomHandler(c *gin.Context)
	DeleteRoomHandler(c *gin.Context)
	GetRoomByIDHandler(c *gin.Context)
}

type GuestProvider interface {
	FetchAllGuests(c *gin.Context)
	CreateGuest(c *gin.Context)
	DeleteGuest(c *gin.Context)
	GetGuestByID(c *gin.Context)
	UpdateGuest(c *gin.Context)
}