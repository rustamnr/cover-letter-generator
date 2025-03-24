package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет наличие access_token в сессии
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		accessToken := session.Get("access_token")

		// Если токен отсутствует в сессии, пробуем взять из заголовка Authorization
		if accessToken == nil {
			authHeader := c.GetHeader("Authorization")
			const bearerPrefix = "Bearer "

			if strings.HasPrefix(authHeader, bearerPrefix) {
				accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
			}
		}

		if accessToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			c.Abort()
			return
		}

		at := fmt.Sprintf("%v", accessToken)
		c.Set("access_token", at)
		c.Next()
	}
}
