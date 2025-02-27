package handlers

import (
	"apigateway/internal/models"
	"apigateway/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	usecase usecase.RoomUsecaseInteface
	logger *slog.Logger
}

func NewRoomHandler(usecase usecase.RoomUsecaseInteface, logger *slog.Logger) RoomProvider {
	return &RoomHandler{usecase: usecase, logger: logger}
}

func (h *RoomHandler) GetRoomsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	rooms, err := h.usecase.GetRooms(ctx)
	if err != nil {
		h.logger.Error("failed to fetch all rooms", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
		return
	}
	h.logger.Info("Succesfully fetched all rooms")
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) CreateRoomHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var room models.RoomDB
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.usecase.CreateRoom(ctx, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *RoomHandler) UpdateRoomHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room ID"})
		return
	}

	var room models.RoomDB
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.usecase.UpdateRoom(ctx, id, room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update room"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "room updated succesfully"})
}

func (h *RoomHandler) DeleteRoomHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room ID"})
		return
	}

	if err := h.usecase.DeleteRoom(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "room deleted succesfully"})
}

func (h *RoomHandler) GetRoomByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing room ID"})
		return
	}
	
	room, err := h.usecase.GetRoomByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch room"})
		return
	}

	if room == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}
