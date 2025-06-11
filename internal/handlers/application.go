package handlers

import (
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (ap *ApplicationHandler) GenerateCoverLetter(c *gin.Context) {
	session := sessions.Default(c)
	ap.service.VacancyProvider.SetAccessToken(session.Get(constants.AccessToken).(string))

	// Get current user resume
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
	resume, err := ap.service.VacancyProvider.GetShortResumeByID(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resume"})
		return
	}

	firstSimilarVacancy, err := ap.service.VacancyProvider.GetFirstShortSuitableVacancy(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	vacancy, err := ap.service.VacancyProvider.GetShortVacancyByID(firstSimilarVacancy.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	coverLetter, err := ap.service.TextGenerator.GenerateCoverLetter(resume, vacancy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating cover letter"})
		return
	}

	c.Set("cover_letter", coverLetter)

	c.JSON(http.StatusOK, gin.H{"cover_letter": coverLetter, "vacancy": vacancy.ID, "resume": resume.ID})
}

func (ap *ApplicationHandler) ApplyToVacancy(c *gin.Context) {
	var (
		err         error
		coverLetter string
		vacancy     *models.VacancyShort
		session     = sessions.Default(c)
	)

	// Set access token from session
	ap.service.VacancyProvider.SetAccessToken(session.Get(constants.AccessToken).(string))

	// Get vacancy by ID from job portal
	vacancyID := c.Param("vacancy_id")
	if vacancyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vacancy ID is required"})
		return
	}
	vacancy, err = ap.service.VacancyProvider.GetShortVacancyByID(vacancyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting vacancy"})
		return
	}
	if vacancy == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vacancy not found"})
		return
	}
	if vacancy.Test.Required {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vacancy requires a test, cannot apply directly"})
		return
	}

	// Get current user resume from session
	currentResume := session.Get(constants.CurrentResumeID)
	if currentResume == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}
	resumeID, ok := currentResume.(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	// Generate cover letter if required
	if vacancy.ResponseLetterRequired {
		// Get resume by ID from job portal
		resume, err := ap.service.VacancyProvider.GetShortResumeByID(resumeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error getting resume": err.Error()})
			return
		}
		if resume == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
			return
		}

		// Generate cover letter using LLM service
		coverLetter, err = ap.service.TextGenerator.GenerateCoverLetter(resume, vacancy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error generating cover letter": err.Error()})
			return
		}
	}

	err = ap.service.VacancyProvider.ApplyToVacancy(resumeID, vacancyID, coverLetter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error applying to vacancy": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "successfully applied to vacancy",
		"vacancy_id": vacancyID,
		"resume_id":  resumeID,
	})
}

func (ap *ApplicationHandler) ApplyToVacansies(c *gin.Context) {
	var (
		// err     error
		session = sessions.Default(c)
	)

	ap.service.VacancyProvider.SetAccessToken(session.Get(constants.AccessToken).(string))

	// Get current user resume from session
	currentResume := session.Get(constants.CurrentResumeID)
	if currentResume == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get user resumes error"})
		return
	}
	resumeID, ok := currentResume.(string)
	if !ok || resumeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume ID not found in session"})
		return
	}

	vacancies, err := ap.service.VacancyProvider.GetShortSimilarVacancies(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error getting similar vacancies": err.Error()})
		return
	}

	for _, vacancy := range vacancies {
		err = ap.service.VacancyProvider.ApplyToVacancy(resumeID, vacancy.ID, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error applying to vacancy": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully applied to all vacancies"})
}
