//test: {
  //.  "email": "testemail01@mail.kz",
  //.   "password": "iitukz"
//}


package http

import (
	"tmp/internal/models"
	"tmp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(router *gin.Engine){
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}




// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.SignUp true "Данные регистрации"
// @Success      201 {object} models.User
// @Failure      400 {object} map[string]string
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.SignUp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body models.SignIn true "SignIn Data"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.SignIn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}