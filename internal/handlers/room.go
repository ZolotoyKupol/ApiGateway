package handlers

import (
	"apigateway/internal/apperr"
	"apigateway/internal/models"
	"apigateway/internal/usecase"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	usecase usecase.RoomUCProvider
	logger  *slog.Logger
}

func NewRoomHandler(usecase usecase.RoomUCProvider, logger *slog.Logger) *RoomHandler {
	return &RoomHandler{usecase: usecase, logger: logger}
}

func (h *RoomHandler) GetAllRooms(c *gin.Context) {
	rooms, err := h.usecase.GetRooms(c)
	if err != nil {
		if errors.Is(err, apperr.ErrNoData) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no rooms found"})
			return
		}
		h.logger.Error("failed to fetch all rooms", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
		return
	}
	h.logger.Debug("Succesfully fetched all rooms")
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var room models.RoomDB
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.usecase.CreateRoom(c, room)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}
	h.logger.Debug("room created successfully", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	var room models.RoomDB
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.usecase.UpdateRoom(c, id, room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update room"})
		return
	}
	h.logger.Debug("room updated successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "room updated succesfully"})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	if err := h.usecase.DeleteRoom(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete room"})
		return
	}
	h.logger.Debug("room deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "room deleted succesfully"})
}

func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	room, err := h.usecase.GetRoomByID(c, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNoData) {
			c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch room"})
		return
	}
	h.logger.Debug("room fetched successfully", "id", id)
	c.JSON(http.StatusOK, room)
}
