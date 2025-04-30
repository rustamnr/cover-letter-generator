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

// SendDeepseekRequest отправляет запрос к DeepSeek API
func (s *DeepSeekService) SendDeepseekRequest(prompt string) (json.RawMessage, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"model": "deepseek-chat",
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": "", // задаёт контекст ассистента (определяет "личность" или роль модели)
				},
				{
					"role":    "assistant",
					"content": "", // предыдущие ответы модели
				},
				{
					"role":    "user",
					"content": prompt, // сообщения от пользователя
				},
			},
			"max_tokens":        20, // between 1 and 8192, default=4096 
			"frequency_penalty": 0,
			"presence_penalty":  0,
			"response_format": map[string]string{
				"type": "text",
			},
			"stop":           nil,
			"stream":         false,
			"stream_options": nil,
			"temperature":    1,
			"top_p":          1,
			"tools":          nil,
			"tool_choice":    "none",
			"logprobs":       false,
			"top_logprobs":   nil,
		}).Post(s.apiURL)

	if err != nil {
		return nil, errors.New("ошибка запроса к DeepSeek")
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("ошибка от DeepSeek API: " + resp.String())
	}

	return json.RawMessage(resp.Body()), nil
}
