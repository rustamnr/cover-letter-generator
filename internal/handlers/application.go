package handlers

import (
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/constants"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (ap *ApplicationHandler) GenerateCoverLetter(c *gin.Context) {
	session := sessions.Default(c)
	ap.service.VacancyProvider.SetAccessToken(session.Get(constants.AccessToken).(string))

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
	resume, err := ap.service.VacancyProvider.GetResume(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting resume"})
		return
	}
	resumePromt := resume.ResumeToLLMModel()

	firstSimilarVacancy, err := ap.service.VacancyProvider.GetFirstSimilarVacancy(resumeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	vacancy, err := ap.service.VacancyProvider.GetVacancyByID(firstSimilarVacancy.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting similar vacancies"})
		return
	}

	vacancyPromt := vacancy.VacancyToLLMModel()

	coverLetter, err := ap.service.TextGenerator.GenerateCoverLetter(resumePromt, vacancyPromt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating cover letter"})
		return
	}

	c.Set("cover_letter", coverLetter)

	c.JSON(http.StatusOK, gin.H{"cover_letter": coverLetter, "vacancy": vacancy.ID, "resume": resume.ID})
}
