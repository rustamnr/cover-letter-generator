package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rustamnr/cover-letter-generator/internal/clients"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/models"
	"github.com/rustamnr/cover-letter-generator/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ApplicationHandler обрабатывает запросы, связанные с заявками
type ApplicationHandler struct {
	service *services.ApplicationService
}

// NewApplicationHandler создает новый ApplicationHandler
func NewApplicationHandler(service *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

// HHHandler handles requests related to hh.ru
type HHHandler struct {
	hhClient *clients.HHClient
}

// NewHHHandler создает новый HHHandler
func NewHHHandler(hhClient *clients.HHClient) *HHHandler {
	return &HHHandler{hhClient: hhClient}
}

// AuthHandler redirects user to the authorization page
func (h *HHHandler) AuthHandler(c *gin.Context) {
	authURL := fmt.Sprintf("https://hh.ru/oauth/authorize?response_type=code&client_id=%s", h.hhClient.ClientID)
	c.Redirect(http.StatusFound, authURL)
}

// CallbackHandler handles the OAuth callback
func (h *HHHandler) CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code not found"})
		return
	}

	accessToken, err := h.hhClient.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.hhClient.GetUserID(accessToken)
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

	h.hhClient.SetAccessToken(accessToken)

	c.JSON(http.StatusOK, gin.H{"message": "authorized", "user_id": userID, "access_token": accessToken})
}

// GetUserResumes retrieves user resumes
func (h *HHHandler) GetUserResumes(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	resumes, err := h.hhClient.GetResumes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resumes"})
		return
	}

	if len(resumes.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "resumes not found"})
		return
	}

	var userResumes []models.SessionResume
	for _, resume := range resumes.Items {
		userResumes = append(userResumes, models.SessionResume{ID: resume.ID, Title: resume.Title})
	}

	session.Set(constants.UserResume, userResumes)
	if err = session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userResumes)
}

// GetCurrentResume retrieves the current resume from session
func (h *HHHandler) GetCurrentResume(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	resumeID, ok := session.Get(constants.CurrentResumeID).(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	resume, err := h.hhClient.GetResume(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"current_resume": resume})
}

// GetVacancyByID retrieves a vacancy by ID
func (h *HHHandler) GetVacancyByID(c *gin.Context) {
	vacancyID := c.Param("vacancy_id")
	if vacancyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vacancy ID is required"})
		return
	}

	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	vacancy, err := h.hhClient.GetVacancyByID(vacancyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vacancy)
}

// SetCurrnetResume select resume by title and save it in session
func (h *HHHandler) SetCurrnetResume(c *gin.Context) {
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

// // GetUserApplications получает список вакансий, на которые пользователь откликнулся
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

// // GetFirstSimilarVacancy get a first similar vacancy
func (h *HHHandler) GetFirstSimilarVacancy(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	resumeID, ok := session.Get(constants.CurrentResumeID).(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	applicationsResponse, err := h.hhClient.GetSuitableVacancies(resumeID, map[string]string{"per_page": "1"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения откликов"})
		return
	}

	nid := applicationsResponse[0].ID
	session.Set("nid", nid)
	session.Save()

	c.JSON(http.StatusOK, applicationsResponse[0])
}

// GetSimilarVacancies get all similar vacancies
func (h *HHHandler) GetSimilarVacancies(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	resumeID, ok := session.Get(constants.CurrentResumeID).(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	vacancies, err := h.hhClient.GetSuitableVacancies(resumeID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}
	if len(vacancies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "similar vacancies not found"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

func (h *HHHandler) CreateCoverLetter(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	currentResume := session.Get(constants.CurrentResumeID)
	if currentResume == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	resumeID, ok := currentResume.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}

	firstSimilarVacancy, err := h.hhClient.GetFirstSuitableVacancy(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	vacancy, err := h.hhClient.GetVacancyByID(firstSimilarVacancy.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	vacancyPromt := vacancy.VacancyToShort()
	_ = vacancyPromt

}
