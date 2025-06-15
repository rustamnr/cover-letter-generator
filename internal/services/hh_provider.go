package services

import (
	clients "github.com/rustamnr/cover-letter-generator/internal/clients/hh"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

type HHProvider struct {
	client *clients.HHClient
}

func NewHHProvider(client *clients.HHClient) *HHProvider {
	return &HHProvider{client: client}
}

func (h *HHProvider) ApplyToVacancy(resumeID, vacancyID, coverLetter string) error {
	return h.client.PostNegotiationByVacancyID(resumeID, vacancyID, coverLetter)
}

func (h *HHProvider) GetResumeByID(resumeID string) (*models.Resume, error) {
	return h.client.GetResume(resumeID)
}

func (h *HHProvider) GetVacancyByID(vacancyID string) (*models.Vacancy, error) {
	return h.client.GetVacancyByID(vacancyID)
}

func (h *HHProvider) GetFirstSuitableVacancy(resumeID string) (*models.Vacancy, error) {
	return h.client.GetFirstSuitableVacancy(resumeID)
}

func (h *HHProvider) GetSimilarVacancies(resumeID, vacanciesLimit string) ([]models.Vacancy, error) {
	return h.client.GetSimilarVacancies(resumeID, map[string]string{
		"per_page": vacanciesLimit,
	})
}

func (h *HHProvider) SetAccessToken(token string) {
	h.client.SetAccessToken(token)
}
