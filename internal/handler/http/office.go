package http

import (
	"net/http"
	"tmp/internal/models"
	"tmp/internal/service"
	"strconv"
	"github.com/gin-gonic/gin"
	"tmp/internal/middleware"
)



type OfficeHandler struct{
	service  *service.OfficeService
}

func NewOfficeHandler(s *service.OfficeService) *OfficeHandler{
	return &OfficeHandler{service: s}
}

func (h *OfficeHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/offices", h.SearchOffices)
	router.GET("/offices/:id", h.GetOffice)
}

func (h *OfficeHandler) SearchOffices(c *gin.Context) {
	var filter models.OfficeFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offices, err := h.service.SearchByAddress(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(offices) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no coworkings found", "offices": []models.Office{}})
		return
	}

	c.JSON(http.StatusOK, offices)
}

func (h *OfficeHandler) GetOffice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid office id"})
		return
	}
	isAdmin := middleware.IsAdmin(c)
	office, err := h.service.GetOffice(c.Request.Context(), id, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "office not found"})
		return
	}
	c.JSON(http.StatusOK, office)
}