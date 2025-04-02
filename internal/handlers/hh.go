package handlers

import (
	"encoding/gob"
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

type SessionResume struct {
	ID    string `json:"id"`
	Title string `json:"title"`
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

// GetUserResumes получает список резюме текущего пользователя
func (h *HHHandler) GetUserResumes(c *gin.Context) {
	// session := sessions.Default(c)
	// accessToken, ok := session.Get(constants.AccessToken).(string)
	// if accessToken == "" {
	// 	authHeader := c.GetHeader("Authorization")
	// 	const bearerPrefix = "Bearer "

	// 	if strings.HasPrefix(authHeader, bearerPrefix) {
	// 		accessToken = strings.TrimPrefix(authHeader, bearerPrefix)
	// 		ok = true
	// 	}
	// }
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует access_token"})
	// 	return
	// }

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
	var userResumes []SessionResume
	for _, resume := range resumesResp {
		userResumes = append(userResumes, SessionResume{ID: resume.ID, Title: resume.Title})
	}

	session.Set("user_resumes", userResumes)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResumes)
}

func init() {
	gob.Register([]SessionResume{})
}

func (h *HHHandler) SelectResume(c *gin.Context) {
	session := sessions.Default(c)

	type titleReq struct {
		Title string `json:"title"`
	}
	var req titleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed parsing request"})
		return
	}
	req.Title = strings.ToLower(req.Title)

	resumesRaw := session.Get("user_resumes")
	if resumesRaw == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	userResumes, ok := resumesRaw.([]SessionResume)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}
	_ = userResumes
	var resumeID string
	for _, resume := range userResumes {
		if strings.Contains(strings.ToLower(resume.Title), req.Title) {
			resumeID = resume.ID
			break
		}
	}
	session.Set("resume_id", resumeID)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving session"})
		return
	}
	c.Set("resume_id", resumeID)
	c.JSON(http.StatusOK, gin.H{"message": "Resume selected", "id": resumeID})

	// c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
}

func (h *HHHandler) GetCurrentResume(c *gin.Context) {
	session := sessions.Default(c)
	accessToken, ok := session.Get(constants.AccessToken).(string)
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

	resumeID, ok := session.Get("resume_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in context"})
		return
	}
	if resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not set in session"})
		return
	}
	resp, err := h.hhService.GetClient().R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(h.hhService.GetAPIURL() + fmt.Sprintf(constants.Resume, resumeID))
	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"error": "error getting resumes"})
		return
	}

	session.Set("current_resume_id", resumeID)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"current resume ID": resumeID})
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
