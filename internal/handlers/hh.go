package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/services"
)

// HHHandler отвечает за обработку запросов к hh.ru API
type HHHandler struct {
	hhService *services.HHService
}

// NewHHHandler создает новый обработчик
func NewHHHandler(hhService *services.HHService) *HHHandler {
	return &HHHandler{hhService: hhService}
}

// AuthHandler редиректит пользователя на hh.ru для авторизации
func (h *HHHandler) AuthHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, h.hhService.GetAuthURL())
}

// CallbackHandler обрабатывает редирект после авторизации
func (h *HHHandler) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует код авторизации"})
		return
	}

	accessToken, err := h.hhService.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.hhService.GetUserID(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения user_id"})
		return
	}

	session := sessions.Default(c)
	session.Set("access_token", accessToken)
	session.Set("user_id", userID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация", "user_id": userID})
}

// GetResumesHandler получает список резюме текущего пользователя
func (h *HHHandler) GetResumesHandler(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + "/resumes/mine")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения резюме"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp.Body())
}
