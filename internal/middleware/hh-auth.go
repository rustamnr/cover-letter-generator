package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
)

// AuthMiddleware проверяет наличие access_token в сессии
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        session := sessions.Default(c)
        accessToken := session.Get(constants.AccessToken)

        // Если токен отсутствует в сессии, пробуем взять из заголовка Authorization
        if accessToken == nil {
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

        if accessToken == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is missing"})
            c.Abort()
            return
        }

        at := fmt.Sprintf("%v", accessToken)
        c.Set(constants.AccessToken, at)
        c.Next()
    }
}
