package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tmp/internal/logger"
	"tmp/internal/middleware"
	"tmp/internal/models"
	"tmp/internal/service"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(s *service.BookingService) *BookingHandler {
	return &BookingHandler{service: s}
}

func (h *BookingHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/bookings", h.Bookings)
	router.POST("/create_booking", h.CreateBooking)
	router.DELETE("/bookings/:id", h.CancelBooking)
}

func (h *BookingHandler) Bookings(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	userID, ok := middleware.GetUserID(c)
	if !ok {
		log.Warn("booking list requested without auth")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookings, err := h.service.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to fetch user bookings", slog.Int("user_id", userID), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "booking not found"})
		return
	}

	log.Debug("user bookings returned", slog.Int("user_id", userID), slog.Int("count", len(bookings)))
	c.JSON(http.StatusOK, bookings)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	userID, ok := middleware.GetUserID(c)
	if !ok {
		log.Warn("booking creation requested without auth")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input models.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("invalid booking payload", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		log.Warn("booking validation failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.CreateBooking(c.Request.Context(), userID, input.RoomID, input.StartTime, input.EndTime, input.DeskID)
	if err != nil {
		if errors.Is(err, service.ErrRoomBooked) {
			log.Warn("booking conflict", slog.Int("room_id", input.RoomID), slog.Int("user_id", userID))
			c.JSON(http.StatusConflict, gin.H{"error": "place already booked"})
			return
		}
		log.Error("failed to create booking", slog.Int("room_id", input.RoomID), slog.Int("user_id", userID), slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("booking created via API", slog.Int("booking_id", booking.ID), slog.Int("user_id", userID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully created",
		"booking": booking,
	})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	userID, ok := middleware.GetUserID(c)
	if !ok {
		log.Warn("booking cancel requested without auth")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Warn("invalid booking id in cancel request", slog.String("value", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := h.service.CancelBooking(c.Request.Context(), id, userID); err != nil {
		log.Warn("cancel booking failed", slog.Int("booking_id", id), slog.Int("user_id", userID), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found or not yours"})
		return
	}

	log.Info("booking cancelled", slog.Int("booking_id", id), slog.Int("user_id", userID))
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}