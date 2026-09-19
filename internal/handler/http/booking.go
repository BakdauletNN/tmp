package http

import (
	"errors"
	"net/http"
	"tmp/internal/middleware"
	"tmp/internal/models"
	"tmp/internal/service"
	
	"github.com/gin-gonic/gin"
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
}

func (h *BookingHandler) Bookings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bookings, err := h.service.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}


func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input models.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.CreateBooking(c.Request.Context(), userID, input.RoomID, input.StartTime, input.EndTime, input.DeskID)
	if err != nil {
		if errors.Is(err, service.ErrRoomBooked) {
			c.JSON(http.StatusConflict, gin.H{"error": "place already booked"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully created",
		"booking": booking,
	})
}

