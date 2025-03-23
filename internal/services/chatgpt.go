package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ChatGPTService отвечает за взаимодействие с API ChatGPT
type ChatGPTService struct {
	apiURL string
	apiKey string
	client *resty.Client
}

// NewChatGPTService создает новый экземпляр ChatGPTService
func NewChatGPTService(apiURL, apiKey string) *ChatGPTService {
	return &ChatGPTService{
		apiURL: apiURL,
		apiKey: apiKey,
		client: resty.New(),
	}
}

// SendRequest отправляет запрос к ChatGPT API
func (s *ChatGPTService) SendRequest(prompt string) (json.RawMessage, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-4",
			"messages":   []map[string]string{{"role": "user", "content": prompt}},
			"max_tokens": 150,
		}).
		Post(s.apiURL)

	if err != nil {
		return nil, errors.New("ошибка запроса к ChatGPT")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("ошибка от ChatGPT API: " + resp.String())
	}

	return json.RawMessage(resp.Body()), nil
}
