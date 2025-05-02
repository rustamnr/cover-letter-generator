package services

import (
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

// VacancyProvider определяет методы для работы с агрегаторами вакансий
type VacancyProvider interface {
	GetResumes(accessToken string) (*models.APIResumeResponse, error)
	GetResume(accessToken, resumeID string) (*models.Resume, error)
	GetVacancyByID(accessToken, vacancyID string) (*models.Vacancy, error)
}

// TextGenerator определяет методы для работы с генераторами текста
type TextGenerator interface {
	GenerateCoverLetter(prompt string) (string, error)
}

// ApplicationService объединяет работу с вакансиями и генерацией текста
type ApplicationService struct {
	VacancyProvider VacancyProvider
	TextGenerator   TextGenerator
}

// NewApplicationService создает новый ApplicationService
func NewApplicationService(vacancyProvider VacancyProvider, textGenerator TextGenerator) *ApplicationService {
	return &ApplicationService{
		VacancyProvider: vacancyProvider,
		TextGenerator:   textGenerator,
	}
}

// ProcessApplication обрабатывает полный цикл: резюме -> вакансия -> сопроводительное письмо -> отклик
func (s *ApplicationService) ProcessApplication(accessToken, resumeID, vacancyID string) error {
	// Получение резюме
	resume, err := s.VacancyProvider.GetResume(accessToken, resumeID)
	if err != nil {
		logger.Errorf("Error getting resume: %v", err)
		return err
	}

	// Получение вакансии
	vacancy, err := s.VacancyProvider.GetVacancyByID(accessToken, vacancyID)
	if err != nil {
		logger.Errorf("Error getting vacancy: %v", err)
		return err
	}

	// Генерация сопроводительного письма
	prompt := "Generate a cover letter for the following resume and vacancy:\n" +
		"Resume: " + resume.Title + "\n" +
		"Vacancy: " + vacancy.Name
	coverLetter, err := s.TextGenerator.GenerateCoverLetter(prompt)
	if err != nil {
		logger.Errorf("Error generating cover letter: %v", err)
		return err
	}

	logger.Infof("Generated cover letter: %s", coverLetter)
	// Здесь можно добавить логику отправки отклика на вакансию

	return nil
}
