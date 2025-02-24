package handlers

import (
	"apigateway/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	usecase *usecase.RoomUsecase
	logger *slog.Logger
}

func NewRoomHandler(usecase *usecase.RoomUsecase, logger *slog.Logger) *RoomHandler {
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
	
}