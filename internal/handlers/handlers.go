package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"apigateway/internal/models"
	"apigateway/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)


type Handlers struct {
	guestUsecase *usecase.GuestUsecase
	logger *slog.Logger
}

func NewHandlers(guestUsecase *usecase.GuestUsecase, logger *slog.Logger) *Handlers {
	return &Handlers{guestUsecase: guestUsecase, logger: logger}
}

func (h *Handlers) FetchAllGuests(c *gin.Context) {
	guests, err := h.guestUsecase.FetchAllGuests()
	if err != nil {
		h.logger.Error("Failed to fetch guests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch guests"})
        return
	}
	h.logger.Info("Successfully handled request to fetch all guests")
	c.JSON(http.StatusOK, guests)
}

func (h *Handlers) CreateGuest(c *gin.Context) {
	var guest models.Guest
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Warn("Invalid input received", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	
	id, err := h.guestUsecase.CreateGuest(guest)
	if err != nil {
		h.logger.Error("Failed to create guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest"})
		return
	}
	h.logger.Info("Guest created successfully", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handlers) DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	err := h.guestUsecase.DeleteGuest(id)
	if err != nil {
		h.logger.Error("Failed to delete guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete guest"})
		return
	}
	h.logger.Info("Guest deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest deleted"})
}

func (h *Handlers) UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var guest models.Guest
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Warn("Invalid input received", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
	}
	err := h.guestUsecase.UpdateGuest(id, guest)
	if err != nil {
		h.logger.Error("Failed to update guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update guest"})
        return
	}
	h.logger.Info("Guest updated successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest updated"})
}

func (h *Handlers) GetGuestByID(c *gin.Context) {
	id := c.Param("id")
	guest, err := h.guestUsecase.GetGuestByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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