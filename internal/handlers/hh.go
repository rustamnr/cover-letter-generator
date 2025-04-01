package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not found"})
		return
	}

	accessToken, err := h.hhService.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.hhService.GetUserID(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting user ID"})
		return
	}

	session := sessions.Default(c)
	session.Set("access_token", accessToken)
	session.Set("user_id", userID)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "authorized", "user_id": userID, "access_token": accessToken})
}

// GetUserResumes получает список резюме текущего пользователя
func (h *HHHandler) GetUserResumes(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		const bearerPrefix = "Bearer "

		if strings.HasPrefix(authHeader, bearerPrefix) {
			accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
			ok = true
		}
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + constants.ResumesMine)
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resumes"})
		return
	}

	var apiResponse models.APIResumeResponse
	if err := json.Unmarshal(resp.Body(), &apiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing HH API response"})
		return
	}

	if len(apiResponse.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resumes not found"})
		return
	}

	resumesResp := apiResponse.Items

	userResumes := make([]models.HHResume, len(resumesResp))
	for i, resume := range resumesResp {
		userResumes[i].ID = resume.ID
		userResumes[i].Title = resume.Title
	}

	c.Set("user_resumes", userResumes)

	c.JSON(http.StatusOK, userResumes)
}

func (h *HHHandler) SelectResume(c *gin.Context) {
	type titleReq struct {
		Title string `json:"title"`
	}
	var req titleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parsing request"})
		return
	}
	req.Title = strings.ToLower(req.Title)

	h.GetUserResumes(c)
	var userResumes []models.HHResume
	resumesRaw, ok := c.Get("user_resumes")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	userResumes, ok = resumesRaw.([]models.HHResume)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	for _, resume := range userResumes {
		if strings.Contains(strings.ToLower(resume.Title), req.Title) {
			session := sessions.Default(c)
			session.Set("resume_id", resume.ID)
			err:=session.Save()
			if err!=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving session"})
				return
			}
			c.Set("resume_id", resume.ID)
			c.JSON(http.StatusOK, gin.H{"message": "Resume selected: " + resume.Title, "id": resume.ID})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
}

func (h *HHHandler) GetCurrentResume(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		const bearerPrefix = "Bearer "

		if strings.HasPrefix(authHeader, bearerPrefix) {
			accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
			ok = true
		}
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	rawId, exist := c.Get("resume_id") // session.Get("resume_id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in context"})
		return
	}
	if rawId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in session"})
		return
	}
	id, ok := rawId.(string)
	if !ok || id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to cast resume ID"})
		return
	}
	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + constants.Resume + "/" + id)
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resumes"})
		return
	}

}

// GetUserApplications получает список вакансий, на которые пользователь откликнулся
func (h *HHHandler) GetUserApplications(c *gin.Context) {
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

// GetUserFirstApplication получает первую вакансию из списка откликов
func (h *HHHandler) GetUserFirstApplication(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		const bearerPrefix = "Bearer "

		if strings.HasPrefix(authHeader, bearerPrefix) {
			accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
			ok = true
		}
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + "/negotiations" + "?" + "per_page=1")

	logger.Infof("resp: %v", resp)

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

	nid := applicationsResponse.Items[0].ID
	session.Set("nid", nid)
	session.Save()

	c.JSON(http.StatusOK, applicationsResponse.Items[0])
}

func (h *HHHandler) SendNewMessage(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get("access_token").(string)
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		const bearerPrefix = "Bearer "

		if strings.HasPrefix(authHeader, bearerPrefix) {
			accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
			ok = true
		}
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
		return
	}

	// Получаем идентификатор отклика из параметров запроса
	// nid := c.Param("nid")
	nid := session.Get("nid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор отклика"})
		return
	}

	if nid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор отклика"})
		return
	}

	var request struct {
		Message string `json:"message"` // Сообщение пользователя
	}

	request.Message = "text"

	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
	// 	return
	// }

	url := fmt.Sprintf(h.hhService.GetAPIURL()+constants.NegotiationsNidMessage, nid)
	// Отправляем сообщение через hhService
	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(map[string]string{"message": request.Message}).
		Post(url)

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки сообщения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Сообщение отправлено"})
}
