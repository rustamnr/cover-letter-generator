package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/models"
	"github.com/rustamnr/cover-letter-generator/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// HHHandler handles requests related
type HHHandler struct {
	hhService *services.HHService
}

// NewHHHandler creates a new HHHandler
func NewHHHandler(hhService *services.HHService) *HHHandler {
	return &HHHandler{hhService: hhService}
}

// AuthHandler rederects user to the authorization page
func (h *HHHandler) AuthHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, h.hhService.GetAuthURL())
}

// CallbackHandler redirects user to the main page after authorization
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set(constants.AccessToken, accessToken)
	session.Set(constants.UserId, userID)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "authorized", "user_id": userID, "access_token": accessToken})
}

// GetUserResumes get user resumes
func (h *HHHandler) GetUserResumes(c *gin.Context) {
	accessToken, exists := c.Get(constants.AccessToken)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	resumes, err := h.hhService.GetResumes(accessToken.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resumes"})
		return
	}

	if len(resumes.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resumes not found"})
		return
	}

	resumesResp := resumes.Items
	var userResumes []models.SessionResume
	for _, resume := range resumesResp {
		userResumes = append(userResumes, models.SessionResume{ID: resume.ID, Title: resume.Title})
	}

	session := sessions.Default(c)
	session.Set(constants.UserResume, userResumes)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResumes)
}

// SelectResume select resume by title and save it in session
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

	session := sessions.Default(c)
	resumesRaw := session.Get(constants.UserResume)
	if resumesRaw == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	userResumes, ok := resumesRaw.([]models.SessionResume)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	var resumeID string
	for _, resume := range userResumes {
		if strings.Contains(strings.ToLower(resume.Title), req.Title) {
			resumeID = resume.ID
			session.Set(constants.CurrentResumeID, resumeID)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving session"})
				return
			}
			session.Set("current_resume_id", resumeID)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving session"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":              "resume selected",
				"current_resume_id":    resumeID,
				"current_resume_title": resume.Title,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
}

// GetCurrentResume get current resume from session
func (h *HHHandler) GetCurrentResume(c *gin.Context) {
	accessToken, exists := c.Get(constants.AccessToken)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access token missing"})
		return
	}

	session := sessions.Default(c)
	resumeID, ok := session.Get(constants.CurrentResumeID).(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in context"})
		return
	}
	if resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in session"})
		return
	}

	resume, err := h.hhService.GetResume(accessToken.(string), resumeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error getting resumes"})
		return
	}
	if resume == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"current resume": resume})
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

// GetUserFirstFoundedApplication получает первую вакансию из списка откликов
func (h *HHHandler) GetUserFirstFoundedApplication(c *gin.Context) {
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
