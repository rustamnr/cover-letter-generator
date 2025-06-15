package handlers_hh

import (
	"net/http"

	clients "github.com/rustamnr/cover-letter-generator/internal/clients/hh"
	"github.com/rustamnr/cover-letter-generator/internal/constants"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HHHandler handles requests related to hh.ru
type HHHandler struct {
	hhClient *clients.HHClient
}

// NewHHHandler создает новый HHHandler
func NewHHHandler(hhClient *clients.HHClient) *HHHandler {
	return &HHHandler{hhClient: hhClient}
}

// GetUserApplications получает список вакансий, на которые пользователь откликнулся
func (h *HHHandler) GetUserApplications(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	applications, err := h.hhClient.GetUserApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения откликов"})
		return
	}

	if len(applications) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Откликов не найдено"})
		return
	}

	c.JSON(http.StatusOK, applications)

	// resp, err := h.hhService.GetClient().R().
	// 	SetHeader("Authorization", "Bearer "+accessToken).
	// 	Get(h.hhService.GetAPIURL() + "/negotiations")

	// if err != nil || resp.StatusCode() != http.StatusOK {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения откликов"})
	// 	return
	// }

	// var applicationsResponse models.APIApplicationsResponse
	// if err := json.Unmarshal(resp.Body(), &applicationsResponse); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки ответа HH API"})
	// 	return
	// }

	// if len(applicationsResponse.Items) == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Откликов не найдено"})
	// 	return
	// }

	// c.JSON(http.StatusOK, applicationsResponse)
}
