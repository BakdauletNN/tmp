package http

import (
	"github.com/gin-gonic/gin"
	"tmp/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/hello", Hi)
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}

func Hi(c *gin.Context) {
	c.JSON(200, gin.H{"hi": "success"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	// TODO
}

func (h *AuthHandler) Login(c *gin.Context) {
	// TODO
}
