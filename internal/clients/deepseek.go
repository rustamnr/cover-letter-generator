package clients

import (
	"encoding/json"
	"fmt"
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

// GenerateCoverLetter отправляет запрос на генерацию сопроводительного письма
func (d *DeepSeekClient) SendPromt(prompt string) (string, error) {
	// Формируем тело запроса
	requestBody := map[string]string{
		"prompt": prompt,
	}

	// Выполняем запрос
	resp, err := d.client.R().
		SetHeader("Authorization", "Bearer "+d.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(d.apiURL + "/generate")

	if err != nil {
		return "", fmt.Errorf("ошибка запроса к DeepSeek API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("ошибка ответа от DeepSeek API: %s", resp.String())
	}

	// Разбираем ответ
	var response struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", fmt.Errorf("ошибка разбора ответа DeepSeek API: %w", err)
	}

	return response.Text, nil
}
