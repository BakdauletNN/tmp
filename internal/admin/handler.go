package admin

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"tmp/internal/logger"
	"tmp/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middleware.IsAdmin(c) {
			logger.FromContext(c.Request.Context()).Warn("admin access blocked")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

type AdminHandler struct {
	repo *AdminRepository
}

func NewAdminHandler(repo *AdminRepository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(RequireAdmin())

	adminGroup.GET("/offices", h.AllOffices)
	adminGroup.GET("/bookings", h.AllBookings)
	adminGroup.DELETE("/bookings/:id", h.CancelBooking)
	adminGroup.GET("/users", h.AllUsers)
}

func (h *AdminHandler) CancelBooking(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil || bookingID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}
	if err := h.repo.CancelBooking(c.Request.Context(), bookingID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		logger.FromContext(c.Request.Context()).Error("admin failed to cancel booking", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not cancel booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

func (h *AdminHandler) AllOffices(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	offices, err := h.repo.AllOffices(c.Request.Context())
	if err != nil {
		log.Error("failed to fetch admin offices", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Debug("admin offices fetched", slog.Int("count", len(offices)))
	c.JSON(http.StatusOK, offices)
}

func (h *AdminHandler) AllBookings(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	bookings, err := h.repo.AllBookings(c.Request.Context())
	if err != nil {
		log.Error("failed to fetch admin bookings", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Debug("admin bookings fetched", slog.Int("count", len(bookings)))
	c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) AllUsers(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	users, err := h.repo.AllUsers(c.Request.Context())
	if err != nil {
		log.Error("failed to fetch admin users", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Debug("admin users fetched", slog.Int("count", len(users)))
	response := make([]gin.H, 0, len(users))
	for _, user := range users {
		response = append(response, gin.H{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"role":          user.Who,
			"is_registered": user.IsRegistered,
		})
	}
	c.JSON(http.StatusOK, response)
}
