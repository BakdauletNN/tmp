package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tmp/internal/logger"
	"tmp/internal/middleware"
	"tmp/internal/models"
	"tmp/internal/service"
)

type OfficeHandler struct {
	service *service.OfficeService
}

func NewOfficeHandler(s *service.OfficeService) *OfficeHandler {
	return &OfficeHandler{service: s}
}

func (h *OfficeHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/offices", h.SearchOffices)
	router.GET("/offices/:id", h.GetOffice)
}

func (h *OfficeHandler) SearchOffices(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var filter models.OfficeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Warn("invalid office search filter", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offices, err := h.service.SearchByAddress(c.Request.Context(), filter)
	if err != nil {
		log.Error("office search failed", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(offices) == 0 {
		log.Debug("no offices found for filter")
		c.JSON(http.StatusOK, gin.H{"message": "no coworkings found", "offices": []models.Office{}})
		return
	}

	log.Debug("offices returned", slog.Int("count", len(offices)))
	c.JSON(http.StatusOK, offices)
}

func (h *OfficeHandler) GetOffice(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Warn("invalid office id", slog.String("value", c.Param("id")))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}
	isAdmin := middleware.IsAdmin(c)
	office, err := h.service.GetOffice(c.Request.Context(), id, isAdmin)
	if err != nil {
		log.Warn("office not found", slog.Int("office_id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "office not found"})
		return
	}
	log.Debug("office info returned", slog.Int("office_id", office.ID))
	c.JSON(http.StatusOK, office)
}
