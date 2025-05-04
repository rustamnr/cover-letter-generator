package middleware

import (
	"net/http"
	"strings"

	"github.com/rustamnr/cover-letter-generator/internal/constants"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет наличие accessToken в сессии или заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// Попробуем получить токен из сессии
		accessToken, ok := session.Get(constants.AccessToken).(string)
		if !ok || accessToken == "" {
			// Если токена нет в сессии, попробуем получить его из заголовка Authorization
			authHeader := c.GetHeader("Authorization")
			const bearerPrefix = "Bearer "

			if strings.HasPrefix(authHeader, bearerPrefix) {
				accessToken = strings.TrimPrefix(authHeader, bearerPrefix)

				// Сохраняем токен в сессии
				session.Set(constants.AccessToken, accessToken)
				if err := session.Save(); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
					c.Abort()
					return
				}
			}
		}

		// Если токен отсутствует, возвращаем ошибку
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is missing"})
			c.Abort()
			return
		}

		// Сохраняем токен в контексте Gin
		c.Set(constants.AccessToken, accessToken)
		c.Next()
	}
}
