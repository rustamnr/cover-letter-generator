package services

import (
	"github.com/rustamnr/cover-letter-generator/internal/clients"
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
func (s *DeepSeekService) GenerateCoverLetter(resume, vacancy string) (string, error) {
	request := clients.LLMRequest{
		System:    "system",
		Assistant: "assistant",
		Content:   "",
		MaxTokens: 20,
	}

	return s.client.SendPromt(request)
}
