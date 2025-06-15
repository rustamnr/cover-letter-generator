package services

import (
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

// JobAgregatorProvider определяет методы для работы с агрегаторами вакансий
type JobAgregatorProvider interface {
	ApplyToVacancy(resumeID, vacancyID, coverLetter string) error
	GetResumeByID(resumeID string) (*models.Resume, error)
	GetVacancyByID(vacancyID string) (*models.Vacancy, error)
	GetFirstSuitableVacancy(resumeID string) (*models.Vacancy, error)
	GetSimilarVacancies(resumeID string, limit string) ([]models.Vacancy, error)
	SetAccessToken(token string)
}

// LLMProvider определяет методы для работы с генераторами текста
type LLMProvider interface {
	GenerateCoverLetter(resume *models.Resume, vacancy *models.Vacancy) (string, error)
}

// ApplicationService объединяет работу с вакансиями и генерацией текста
type ApplicationService struct {
	VacancyProvider JobAgregatorProvider
	TextGenerator   LLMProvider
}

// NewApplicationService создает новый ApplicationService
func NewApplicationService(vacancyProvider JobAgregatorProvider, textGenerator LLMProvider) *ApplicationService {
	return &ApplicationService{
		VacancyProvider: vacancyProvider,
		TextGenerator:   textGenerator,
	}
}
