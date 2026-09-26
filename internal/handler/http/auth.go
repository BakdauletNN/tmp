//test: {
//.  "email"testemail01@mail.kz: "",
//.   "password": "iitukz"
//}

package http

import (
	"errors"
	"log/slog"
	"net/http"

	"tmp/internal/logger"
	"tmp/internal/models"
	"tmp/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
	router.POST("/forgot-password", h.ForgotPassword)
	router.POST("/reset-password", h.ResetPassword)
}

type forgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var input forgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid email is required"})
		return
	}

	if err := h.service.RequestPasswordReset(c.Request.Context(), input.Email); err != nil {
		logger.FromContext(c.Request.Context()).Error(
			"password reset email could not be sent",
			slog.String("error", err.Error()),
		)
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, we've sent a link to reset the password."})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input resetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token and password (at least 6 characters) are required"})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), input.Token, input.Password); err != nil {
		if errors.Is(err, service.ErrInvalidPasswordReset) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reset link is invalid or expired"})
			return
		}
		logger.FromContext(c.Request.Context()).Error("password reset failed", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

// Register godoc
// @Summary      Register a user
// @Description  Creates a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.SignUp true "Registration data"
// @Success      201 {object} models.User
// @Failure      400 {object} map[string]string
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var input models.SignUp
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("invalid register payload", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		log.Warn("register validation failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		log.Error("register failed", slog.String("email", input.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("user registered", slog.Int("user_id", user.ID), slog.String("email", user.Email))
	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Register godoc
// @Summary      Register a user
// @Description  Creates a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.SignIn true "SignIn Data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var input models.SignIn
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Warn("invalid login payload", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		log.Warn("login validation failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		log.Warn("login failed", slog.String("email", input.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Info("user logged in", slog.String("email", input.Email))
	c.JSON(http.StatusOK, gin.H{"token": token})
}
