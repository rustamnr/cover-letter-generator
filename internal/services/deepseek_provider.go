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
func (s *DeepSeekService) SendRequest(prompt string) (string, error) {
	s.client.SendPromt(prompt)
	// resp, err := s.client.R().
	// 	SetHeader("Authorization", "Bearer "+s.apiKey).
	// 	SetHeader("Content-Type", "application/json").
	// 	SetBody(map[string]any{
	// 		"model": "deepseek-chat",
	// 		"messages": []map[string]string{
	// 			{
	// 				"role":    "system",
	// 				"content": "", // задаёт контекст ассистента (определяет "личность" или роль модели)
	// 			},
	// 			{
	// 				"role":    "assistant",
	// 				"content": "", // предыдущие ответы модели
	// 			},
	// 			{
	// 				"role":    "user",
	// 				"content": prompt, // сообщения от пользователя
	// 			},
	// 		},
	// 		"max_tokens":        20, // between 1 and 8192, default=4096
	// 		"frequency_penalty": 0,
	// 		"presence_penalty":  0,
	// 		"response_format": map[string]string{
	// 			"type": "text",
	// 		},
	// 		"stop":           nil,
	// 		"stream":         false,
	// 		"stream_options": nil,
	// 		"temperature":    1,
	// 		"top_p":          1,
	// 		"tools":          nil,
	// 		"tool_choice":    "none",
	// 		"logprobs":       false,
	// 		"top_logprobs":   nil,
	// 	}).Post(s.apiURL)

	// if err != nil {
	// 	return nil, errors.New("ошибка запроса к DeepSeek")
	// }

	// if resp.StatusCode() != http.StatusOK {
	// 	return nil, errors.New("ошибка от DeepSeek API: " + resp.String())
	// }

	// return json.RawMessage(resp.Body()), nil

	return "", nil
}
