package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"apigateway/internal/apperr"
	"apigateway/internal/models"
	"apigateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

type GuestHandlers struct {
	guestUsecase usecase.GuestUCProvider
	logger       *slog.Logger
}

func NewHandlers(guestUsecase usecase.GuestUCProvider, logger *slog.Logger) *GuestHandlers {
	return &GuestHandlers{guestUsecase: guestUsecase, logger: logger}
}

func (h *GuestHandlers) GetAllGuests(c *gin.Context) {
	guests, err := h.guestUsecase.GetGuests(c)
	if err != nil {
		if errors.Is(err, apperr.ErrNoData) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no guests found"})
			return
		}
		h.logger.Error("failed to fetch guests", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch guests"})
		return
	}
	h.logger.Info("successfully handled request to fetch all guests")
	c.JSON(http.StatusOK, guests)
}

func (h *GuestHandlers) CreateGuest(c *gin.Context) {
	var guest models.GuestDB
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Debug("invalid input received", "input", guest, "error", err)
		h.logger.Warn("invalid input received", "error", err.Error(), "request body", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "details": err.Error()})
		return
	}
	h.logger.Debug("received guest data", "guest", guest)

	
	id, err := h.guestUsecase.CreateGuest(c, guest)
	if err != nil {
		h.logger.Error("failed to create guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest"})
		return
	}
	h.logger.Debug("guest created successfully", "id", id)
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *GuestHandlers) DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	err := h.guestUsecase.DeleteGuest(c, id)
	if err != nil {
		h.logger.Error("failed to delete guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete guest"})
		return
	}
	h.logger.Debug("guest deleted successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest deleted"})
}

func (h *GuestHandlers) UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var guest models.GuestDB
	if err := c.BindJSON(&guest); err != nil {
		h.logger.Warn("invalid input received", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	err := h.guestUsecase.UpdateGuest(c, id, guest)
	if err != nil {
		h.logger.Error("failed to update guest", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update guest"})
		return
	}
	h.logger.Debug("guest updated successfully", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "guest updated"})
}

func (h *GuestHandlers) GetGuestByID(c *gin.Context) {
	id := c.Param("id")
	guest, err := h.guestUsecase.GetGuestByID(c, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNoData) {
			h.logger.Warn("guest not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "guest not found"})
			return
		}
		h.logger.Error("failed to fetch guest by ID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "failed to fetch guest"})
		return
	}
	h.logger.Debug("successfully fetched guest", "id", id)
	c.JSON(http.StatusOK, guest)
}
