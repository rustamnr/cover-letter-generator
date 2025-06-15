package clients_hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/constants"
)

// ===== Authentication and Token Management =====

func (c *HHClient) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *HHClient) GetAccessToken() string {
	return c.accessToken
}

func (c *HHClient) ExchangeCodeForToken(code string) (string, error) {
	resp, err := c.client.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     c.ClientID,
			"client_secret": c.clientSecret,
			"code":          code,
			"redirect_uri":  c.redirectURI,
		}).
		Post(constants.HHOAuth)

	if err != nil {
		return "", fmt.Errorf("ошибка запроса к hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("не удалось получить токен: %s", resp.String())
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &tokenData); err != nil {
		return "", fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		return "", errors.New("access_token не найден")
	}

	return accessToken, nil
}
