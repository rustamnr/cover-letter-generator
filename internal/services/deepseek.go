package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// DeepSeekService отвечает за взаимодействие с API DeepSeek
type DeepSeekService struct {
	apiURL string
	apiKey string
	client *resty.Client
}

// NewDeepSeekService создает новый экземпляр DeepSeekService
func NewDeepSeekService(apiURL, apiKey string) *DeepSeekService {
	return &DeepSeekService{
		apiURL: apiURL,
		apiKey: apiKey,
		client: resty.New(),
	}
}

// SendRequest отправляет запрос к DeepSeek API
func (s *DeepSeekService) SendRequest(prompt string) (json.RawMessage, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "deepseek-chat",
			"messages":   []map[string]string{{"role": "user", "content": prompt}},
			"max_tokens": 150,
		}).
		Post(s.apiURL)

	if err != nil {
		return nil, errors.New("ошибка запроса к DeepSeek")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("ошибка от DeepSeek API: " + resp.String())
	}

	return json.RawMessage(resp.Body()), nil
}
