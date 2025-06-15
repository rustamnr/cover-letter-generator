package handlers_hh

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
)

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

// // GetFirstSimilarVacancy get a first similar vacancy
func (h *HHHandler) GetFirstSimilarVacancy(c *gin.Context) {
	session := sessions.Default(c)
	h.hhClient.SetAccessToken(session.Get(constants.AccessToken).(string))

	resumeID, ok := session.Get(constants.CurrentResumeID).(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	applicationsResponse, err := h.hhClient.GetSimilarVacancies(resumeID, map[string]string{"per_page": "1"})
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

	perPage := c.Query("per_page") // Default 10, max 100
	if perPage == "" {
		perPage = "100"
	}

	vacancies, err := h.hhClient.GetSimilarVacancies(resumeID,
		map[string]string{"per_page": perPage},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getting similar vacancies": err.Error()})
		return
	}
	if len(vacancies) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "similar vacancies not found"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}
