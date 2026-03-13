package handlers

import (
	"errors"
	"l3/EventBooker/internal/customErrs"
	"l3/EventBooker/internal/models"
	"l3/EventBooker/internal/service"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type EventHandler struct {
	bookingService *service.BookingService
}

func NewEventHandler(bookingService *service.BookingService) *EventHandler {
	return &EventHandler{
		bookingService: bookingService,
	}
}

func (h *EventHandler) GetAllEvents(c *ginext.Context) {
	events, err := h.bookingService.GetAllEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "Failed to get events",
		})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) CreateEvent(c *ginext.Context) {
	var req CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	if req.Title == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "Title is required",
		})
		return
	}
	if req.TotalSeats <= 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "TotalSeats must be positive",
		})
		return
	}
	if req.Date.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "Date must be in the future",
		})
		return
	}
	if req.Price < 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "Price must be positive",
		})
		return
	}
	event := &models.Event{
		Title:           req.Title,
		Description:     req.Description,
		Date:            req.Date,
		TotalSeats:      req.TotalSeats,
		AvailableSeats:  req.TotalSeats,
		Price:           req.Price,
		PaymentRequired: req.PaymentRequired,
		CreatedAt:       time.Now(),
	}

	err := h.bookingService.CreateEvent(c.Request.Context(), event)
	if err != nil {
		//надо ли мне здесь вообще логировать(по идее да)
		zlog.Logger.Error().Err(err).Msg("Failed to create event")
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "Failed to create event",
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Book(c *ginext.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "Username is required",
		})
		return
	}

	booking, err := h.bookingService.Book(c.Request.Context(), id, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error":   "Failed to book",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (h *EventHandler) Confirm(c *ginext.Context) {
	bookingID := c.Param("id")
	if len(bookingID) == 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	err := h.bookingService.Confirm(c.Request.Context(), bookingID)
	if err != nil {
		switch {
		case errors.Is(err, customErrs.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, ginext.H{"error": "Booking not found"})
		case errors.Is(err, customErrs.ErrBookingAlreadyConfirmed):
			c.JSON(http.StatusConflict, ginext.H{"error": "Booking already confirmed"})
		case errors.Is(err, customErrs.ErrBookingExpired):
			c.JSON(http.StatusConflict, ginext.H{"error": "Booking has expired"})
		case errors.Is(err, customErrs.ErrBookingCanceled):
			c.JSON(http.StatusConflict, ginext.H{"error": "Booking was canceled"})
		default:
			c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to confirm booking"})
		}
		return
	}
	c.JSON(http.StatusOK, ginext.H{
		"message": "Booking confirmed successfully",
	})
}

func (h *EventHandler) GetEventWithDetails(c *ginext.Context) {
	id := c.Param("id")
	if len(id) == 0 {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "id is empty",
		})
		return
	}
	event, err := h.bookingService.GetEventWithDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error":   "failed to get event details",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, event)
}
