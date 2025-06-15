package services

import (
	"fmt"

	clients "github.com/rustamnr/cover-letter-generator/internal/clients/deepseek"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/models"
	"github.com/rustamnr/cover-letter-generator/pkg/promts"
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
func (s *DeepSeekService) GenerateCoverLetter(resume *models.Resume, vacancy *models.Vacancy) (string, error) {
	content := fmt.Sprint(resume.ToString(), vacancy.ToString())
	logger.Debugf("Deepseek request content: %s", content)
	request := clients.LLMRequest{
		System:    promts.DeepseekSystemContext,
		Content:   content,
		MaxTokens: 2048,
	}

	return s.client.SendPromt(request)
}
