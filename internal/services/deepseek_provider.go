package services

import (
	"fmt"
	"strings"

	"github.com/rustamnr/cover-letter-generator/internal/clients"
	"github.com/rustamnr/cover-letter-generator/internal/models"
	"github.com/rustamnr/cover-letter-generator/pkg/templates"
)

// DeepSeekService отвечает за взаимодействие с API DeepSeek
type DeepSeekService struct {
	client *clients.DeepSeekClient
}

// NewDeepSeekService создает новый экземпляр DeepSeekService
func NewDeepSeekService(client *clients.DeepSeekClient) *DeepSeekService {
	return &DeepSeekService{
		client: client,
	}
}

// SendDeepseekRequest отправляет запрос к DeepSeek API
func (s *DeepSeekService) GenerateCoverLetter(resume models.ResumeForLLM, vacancy models.VacancyForLLM) (string, error) {
	request := clients.LLMRequest{
		System: templates.BasePromt,

		Content: fmt.Sprintf(
			"ВАКАНСИЯ:\n"+
				"- Название: %s\n"+
				"- Компания: %s\n"+
				"- Локация: %s\n"+
				"- Описание: %s\n"+
				"- Требуемые навыки: %s\n\n"+
				"РЕЗЮМЕ КАНДИДАТА:\n"+
				"- Имя: %s %s\n"+
				"- Локация: %s\n"+
				"- Опыт работы: %d лет\n"+
				"- Навыки: %s\n"+
				vacancy.Name,
			vacancy.CompanyName,
			vacancy.Location,
			vacancy.Description,
			strings.Join(vacancy.KeySkills, ", "),
			resume.FirstName,
			resume.LastName,
			resume.Location,
			resume.TotalExperience/12,
			strings.Join(resume.Skills, ", "),
		),
		MaxTokens: 2048,
	}

	return s.client.SendPromt(request)
}

// limitString ограничивает длину строки до maxLength символов
func limitString(input string, maxLength int) string {
	if len(input) > maxLength {
		return input[:maxLength] + "..."
	}
	return input
}
