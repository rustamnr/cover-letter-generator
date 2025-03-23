package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/models"
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

	var apiResponse models.APIResumeResponse
	if err := json.Unmarshal(resp.Body(), &apiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ответа HH API"})
		return
	}

	if len(apiResponse.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Резюме не найдено"})
		return
	}

	resumeData := apiResponse.Items[0] // Берем первое резюме

	// Обрабатываем контакты
	for i, contact := range resumeData.Contact {
		var phone models.PhoneValue
		var email string

		if err := json.Unmarshal(contact.Value, &phone); err == nil && phone.Number != "" {
			// Телефон
			resumeData.Contact[i].ParsedValue = fmt.Sprintf("+%s (%s) %s", phone.Country, phone.City, phone.Number)
		} else if err := json.Unmarshal(contact.Value, &email); err == nil {
			// Email
			resumeData.Contact[i].ParsedValue = email
		} else {
			resumeData.Contact[i].ParsedValue = "Неизвестный формат"
		}
	}

	c.JSON(http.StatusOK, resumeData) // Отправляем JSON
}

// GetUserApplicationsHandler получает список вакансий, на которые пользователь откликнулся
func (h *HHHandler) GetUserApplicationsHandler(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + "/negotiations")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения откликов"})
		return
	}

	var applicationsResponse models.APIApplicationsResponse
	if err := json.Unmarshal(resp.Body(), &applicationsResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ответа HH API"})
		return
	}

	if len(applicationsResponse.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Откликов не найдено"})
		return
	}

	c.JSON(http.StatusOK, applicationsResponse)
}