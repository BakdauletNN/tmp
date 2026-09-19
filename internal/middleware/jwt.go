package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userCtxKey = "user_id"

// ParseToken checks signature
func ParseToken(tokenStr, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := ParseToken(tokenStr, secret)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set(userCtxKey, int(claims["user_id"].(float64)))
		c.Next()
	}
}

func GetUserID(c *gin.Context) (int, bool) {
	v, ok := c.Get(userCtxKey)
	if !ok {
		return 0, false
	}
	return v.(int), true
}