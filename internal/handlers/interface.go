package handlers

import "github.com/gin-gonic/gin"

type RoomProvider interface {
	GetAllRooms(c *gin.Context)
	CreateRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
	GetRoomByID(c *gin.Context)
}

type GuestProvider interface {
	GetAllGuests(c *gin.Context)
	CreateGuest(c *gin.Context)
	DeleteGuest(c *gin.Context)
	GetGuestByID(c *gin.Context)
	UpdateGuest(c *gin.Context)
}
