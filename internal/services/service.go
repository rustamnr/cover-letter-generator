package services

import (
	"fmt"

	"github.com/rustamnr/cover-letter-generator/internal/models"
)

// VacancyProvider определяет методы для работы с агрегаторами вакансий
type VacancyProvider interface {
	GetResumes(accessToken string) (*models.APIResumeResponse, error)
	GetResume(accessToken, resumeID string) (*models.Resume, error)
	// GetCurrentResume(accessToken string) (*models.Resume, error)
	GetVacancyByID(accessToken, vacancyID string) (*models.Vacancy, error)
	GetSimilarVacancies(accessToken, resumeID string, queryParams map[string]string) (*models.SimilarVacanciesResponse, error)
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
func (s *ApplicationService) GenerateCoverLetter(accessToken, resumeID string) (string, error) {
	// 1. Получить текущее резюме
	resume, err := s.VacancyProvider.GetResume(accessToken, resumeID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения резюме: %w", err)
	}

	// 2. Получить рекомендованную вакансию
	queryParams := map[string]string{
		"per_page": "1",
	}
	similarVacancies, err := s.VacancyProvider.GetSimilarVacancies(accessToken, resumeID, queryParams)
	if err != nil {
		return "", fmt.Errorf("ошибка получения рекомендованных вакансий: %w", err)
	}
	if len(similarVacancies.Items) == 0 {
		return "", fmt.Errorf("рекомендованные вакансии не найдены")
	}
	vacancy := similarVacancies.Items[0]

	// 3. Преобразовать резюме и вакансию в модели для LLM
	resumeForLLM := resume.ResumeToLLMModel()
	vacancyForLLM := vacancy.VacancyToLLMModel()

	// 4. Сформировать запрос для LLM
	prompt := fmt.Sprintf(
		"Составь сопроводительное письмо для вакансии:\n\nВакансия:\n%+v\n\nРезюме:\n%+v",
		vacancyForLLM,
		resumeForLLM,
	)

	// 5. Отправить запрос в DeepSeek
	coverLetter, err := s.TextGenerator.GenerateCoverLetter(prompt)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации сопроводительного письма: %w", err)
	}

	// 6. Вернуть результат
	return coverLetter, nil
}
