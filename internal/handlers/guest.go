package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"apigateway/internal/models"
	"apigateway/internal/usecase"

	"github.com/gin-gonic/gin"
)


type GuestHandlers struct {
	guestUsecase usecase.GuestUsecaseInterface
	logger *slog.Logger
}

func NewHandlers(guestUsecase usecase.GuestUsecaseInterface, logger *slog.Logger) GuestProvider {
	return &GuestHandlers{guestUsecase: guestUsecase, logger: logger}
}

func (h *GuestHandlers) FetchAllGuests(c *gin.Context) {
	ctx := c.Request.Context()
	guests, err := h.guestUsecase.FetchAllGuests(ctx)
	if err != nil {
		h.logger.Error("Failed to fetch guests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch guests"})
        return
	}
	h.logger.Info("Successfully handled request to fetch all guests")
	c.JSON(http.StatusOK, guests)
}

func (h *GuestHandlers) CreateGuest(c *gin.Context) {
	ctx := c.Request.Context()
	var guest models.GuestDB
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Warn("Invalid input received", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	
	id, err := h.guestUsecase.CreateGuest(ctx, guest)
	if err != nil {
		h.logger.Error("Failed to create guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest"})
		return
	}
	h.logger.Info("Guest created successfully", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *GuestHandlers) DeleteGuest(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := h.guestUsecase.DeleteGuest(ctx, id)
	if err != nil {
		h.logger.Error("Failed to delete guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete guest"})
		return
	}
	h.logger.Info("Guest deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest deleted"})
}

func (h *GuestHandlers) UpdateGuest(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var guest models.GuestDB
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Warn("Invalid input received", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
	}
	err := h.guestUsecase.UpdateGuest(ctx, id, guest)
	if err != nil {
		h.logger.Error("Failed to update guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update guest"})
        return
	}
	h.logger.Info("Guest updated successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest updated"})
}

func (h *GuestHandlers) GetGuestByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	guest, err := h.guestUsecase.GetGuestByID(ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrNoData) {
			h.logger.Warn("Guest not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "guest not found"})
            return
		}
		h.logger.Error("Failed to fetch guest by ID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "failed to fetch guest"})
		return
	}
	h.logger.Info("Successfully fetched guest", "id", id)
	c.JSON(http.StatusOK, guest)
}