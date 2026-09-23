package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tmp/internal/middleware"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !middleware.IsAdmin(c) {
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
	adminGroup.GET("/users", h.AllUsers)
}

func (h *AdminHandler) AllOffices(c *gin.Context) {
	offices, err := h.repo.AllOffices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, offices)
}

func (h *AdminHandler) AllBookings(c *gin.Context) {
	bookings, err := h.repo.AllBookings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) AllUsers(c *gin.Context) {
	users, err := h.repo.AllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}