package http

import (
	"net/http"
	"strconv"

	"tmp/internal/models"
	"tmp/internal/service"
	"tmp/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct{
	service *service.RoomService
}

func NewRoomHanlder(s *service.RoomService) *RoomHandler{
	return &RoomHandler{service: s}
}

func (h *RoomHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/rooms", h.SearchRoom)    
	router.GET("/rooms/:id", h.GetRoomInfo)
}

func (h *RoomHandler) SearchRoom(c *gin.Context) {
	var filter models.RoomFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rooms, err := h.service.SearchRoom(c.Request.Context(), &filter)  
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rooms) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no rooms found", "rooms": []models.Room{}}) 
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoomInfo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	room, err := h.service.GetRoom(c.Request.Context(), id, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}


type SetCodeInput struct {
	Code int `json:"code" validate:"required" example:"1234"`
}

func (h *RoomHandler) SetCode(c *gin.Context) {
	idParam := c.Param("id")
	roomID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input SetCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetAccessCode(c.Request.Context(), roomID, input.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "code updated"})
}


