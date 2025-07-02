package middleware

import (
	"auth-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение токена из заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Проверка формата "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Валидация токена
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохранение информации о пользователе в контексте
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}