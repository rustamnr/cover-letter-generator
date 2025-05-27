package services

import (
	"github.com/rustamnr/cover-letter-generator/internal/clients"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

type HHProvider struct {
	client *clients.HHClient
}

func NewHHProvider(client *clients.HHClient) *HHProvider {
	return &HHProvider{client: client}
}

func (h *HHProvider) GetResumeByID(resumeID string) (*models.Resume, error) {
	return h.client.GetResume(resumeID)
}

func (h *HHProvider) GetShortResumeByID(resumeID string) (*models.ResumeShort, error) {
	return h.client.GetShortResume(resumeID)
}

func (h *HHProvider) GetVacancyByID(vacancyID string) (*models.Vacancy, error) {
	return h.client.GetVacancyByID(vacancyID)
}

func (h *HHProvider) GetShortVacancyByID(vacancyID string) (*models.VacancyShort, error) {
	return h.client.GetShortVacancyByID(vacancyID)
}

func (h *HHProvider) GetFirstShortSuitableVacancy(resumeID string) (*models.VacancyShort, error) {
	return h.client.GetFirstShortSuitableVacancy(resumeID)
}

func (h *HHProvider) ApplyToVacancy(resumeID, vacancyID, coverLetter string) error {
	return h.client.PostNegotiationByVacancyID(resumeID, vacancyID, coverLetter)
}

func (h *HHProvider) SetAccessToken(token string) {
	h.client.SetAccessToken(token)
}
