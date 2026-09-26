package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"tmp/internal/logger"
	"tmp/internal/repository"
)

type DeskHandler struct {
	repo *repository.DeskRepository
}

func NewDeskHandler(repo *repository.DeskRepository) *DeskHandler {
	return &DeskHandler{repo: repo}
}

func (h *DeskHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/rooms/:id/desks", h.ListAvailability)
}

func (h *DeskHandler) ListAvailability(c *gin.Context) {
	roomID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roomID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}
	start, err := time.Parse(time.RFC3339, c.Query("start_time"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time must be RFC3339"})
		return
	}
	end, err := time.Parse(time.RFC3339, c.Query("end_time"))
	if err != nil || !end.After(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be after start_time and use RFC3339"})
		return
	}

	desks, err := h.repo.ListAvailability(c.Request.Context(), roomID, start, end)
	if err != nil {
		logger.FromContext(c.Request.Context()).Error("failed to list desk availability")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load desk availability"})
		return
	}
	c.JSON(http.StatusOK, desks)
}
