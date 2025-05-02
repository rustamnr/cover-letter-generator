package clients

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/go-resty/resty/v2"
)

// DeepSeekService представляет клиента для работы с API DeepSeek
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

// GenerateCoverLetter отправляет запрос на генерацию сопроводительного письма
func (d *DeepSeekService) GenerateCoverLetter(prompt string) (string, error) {
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
