package clients

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

// DeepSeekClient представляет клиента для работы с API DeepSeek
type DeepSeekClient struct {
	apiURL string
	apiKey string
	client *resty.Client
}

// NewDeepSeekClient создает новый экземпляр DeepSeekService
func NewDeepSeekClient() *DeepSeekClient {
	return &DeepSeekClient{
		apiURL: os.Getenv("DEEPSEEK_API_URL"),
		apiKey: os.Getenv("DEEPSEEK_API_KEY"),
		client: resty.New(),
	}
}

type LLMRequest struct {
	System    string // задаёт контекст ассистента (определяет "личность" или роль модели)
	Assistant string // предыдущие ответы модели
	Content   string // сообщения от пользователя
	MaxTokens int    // between 1 and 8192, default=4096
}

// GenerateCoverLetter отправляет запрос на генерацию сопроводительного письма
func (d *DeepSeekClient) SendPromt(req LLMRequest) (string, error) {
	resp, err := d.client.R().
		SetHeader("Authorization", "Bearer "+d.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"model": "deepseek-chat",
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": req.System,
				},
				{
					"role":    "assistant",
					"content": req.Assistant,
				},
				{
					"role":    "user",
					"content": req.Content,
				},
			},
			"max_tokens":        req.MaxTokens,
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
		}).Post(d.apiURL)

	if err != nil {
		return "", errors.New("ошибка запроса к DeepSeek")
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("ошибка от DeepSeek API: " + resp.String())
	}

	return "", nil
}
