package handlers_hh

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

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
