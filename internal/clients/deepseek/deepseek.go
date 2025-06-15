package clients_deepseek

import (
	"encoding/json"
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

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
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
				// "type": "json_object",
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

	// Разбираем JSON-ответ
	var response DeepSeekResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", errors.New("ошибка разбора ответа DeepSeek API")
	}

	// Извлекаем текст сопроводительного письма
	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return "", errors.New("ответ от DeepSeek API не содержит текста")
	}

	return response.Choices[0].Message.Content, nil
}
