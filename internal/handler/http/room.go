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

type RoomHandler struct {
	service *service.RoomService
}

func NewRoomHanlder(s *service.RoomService) *RoomHandler {
	return &RoomHandler{service: s}
}

func (h *RoomHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/rooms", h.SearchRoom)
	router.GET("/rooms/:id", h.GetRoomInfo)
}

func (h *RoomHandler) SearchRoom(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var filter models.RoomFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Warn("invalid room search filter", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rooms, err := h.service.SearchRoom(c.Request.Context(), &filter)
	if err != nil {
		log.Error("room search failed", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rooms) == 0 {
		log.Debug("no rooms found for filter")
		c.JSON(http.StatusOK, gin.H{"message": "no rooms found", "rooms": []models.Room{}})
		return
	}

	log.Debug("rooms found", slog.Int("count", len(rooms)))
	c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoomInfo(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid room id", slog.String("value", idParam))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	isAdmin := middleware.IsAdmin(c)
	room, err := h.service.GetRoom(c.Request.Context(), id, isAdmin)
	if err != nil {
		log.Warn("room not found", slog.Int("room_id", id), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	log.Debug("room info returned", slog.Int("room_id", room.ID))
	c.JSON(http.StatusOK, room)
}

type SetCodeInput struct {
	Code int `json:"code" validate:"required" example:"1234"`
}

func (h *RoomHandler) SetCode(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	idParam := c.Param("id")
	roomID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Warn("invalid room id on set code", slog.String("value", idParam))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room id"})
		return
	}

	var input SetCodeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("invalid room code payload", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetAccessCode(c.Request.Context(), roomID, input.Code); err != nil {
		log.Error("failed to update room access code", slog.Int("room_id", roomID), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("room access code updated", slog.Int("room_id", roomID))
	c.JSON(http.StatusOK, gin.H{"message": "code updated"})
}
