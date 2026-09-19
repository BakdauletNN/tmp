package http

import (
	"strconv"
	"net/http"

	"tmp/internal/service"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct{
	service *service.RoomService
}

func NewRoomHanlder(s *service.RoomService) *RoomHandler{
	return &RoomHandler{service: s}
}

func (h *RoomHandler) RegisterRoutes(router *gin.Engine){
	router.GET("/room")
	router.GET("/rooms/:id", h.GetRoomInfo)
}

func (h *RoomHandler) GetRoomInfo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	room, err := h.service.GetRoom(c.Request.Context(), id)
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