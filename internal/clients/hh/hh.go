package clients_hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

type HHClient struct {
	apiURL       string
	ClientID     string
	clientSecret string
	redirectURI  string
	client       *resty.Client
	accessToken  string
}

func NewHHClient() *HHClient {
	return &HHClient{
		apiURL:       os.Getenv("HH_API_URL"),
		ClientID:     os.Getenv("HH_CLIENT_ID"),
		clientSecret: os.Getenv("HH_CLIENT_SECRET"),
		client:       resty.New(),
	}
}

func (c *HHClient) GetUserID(accessToken string) (string, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(c.apiURL + constants.Me)

	if err != nil {
		return "", fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("не удалось получить user_id: %s", resp.String())
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &userData); err != nil {
		return "", fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	userID, ok := userData["id"].(string)
	if !ok {
		return "", errors.New("user_id не найден")
	}

	return userID, nil
}

func (c *HHClient) GetUserApplications() ([]models.ApplicationItem, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		Get(c.apiURL + constants.Negotiations)

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить список заявок: %s", resp.String())
	}

	var applications []models.ApplicationItem
	if err := json.Unmarshal(resp.Body(), &applications); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return applications, nil
}

func (c *HHClient) GetFirstFoundedApplication(accessToken string) (*models.APIApplicationsResponse, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(c.apiURL + "/negotiations" + "?" + "per_page=1")

	logger.Infof("resp: %v", resp)

	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	var applicationsResponse models.APIApplicationsResponse
	if err := json.Unmarshal(resp.Body(), &applicationsResponse); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	if len(applicationsResponse.Items) == 0 {
		return nil, errors.New("заявки не найдены")
	}

	return &applicationsResponse, nil
}

func (c *HHClient) PostNegotiationByVacancyID(resumeID, vacancyID, message string) error {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetMultipartFormData(map[string]string{
			"resume_id":  resumeID,
			"vacancy_id": vacancyID,
			"message":    message,
		}).
		Post(c.apiURL + "/negotiations")
	if err != nil {
		return fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}
	if resp.StatusCode() != http.StatusCreated { // https://api.hh.ru/negotiations has 201 response code on success
		logger.Errorf("не удалось создать заявку: %s", resp.String())
		return fmt.Errorf("не удалось создать заявку: %s", resp.String())
	}

	return nil
}

func (c *HHClient) SendMessage() {

}
