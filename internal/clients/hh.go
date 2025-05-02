package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
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
	access_token string
}

func init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("ошибка загрузки .env файла: %v", err))
	}
}

// NewHHClient создает новый экземпляр HHClient
func NewHHClient() *HHClient {
	return &HHClient{
		apiURL:       os.Getenv("HH_API_URL"),
		ClientID:     os.Getenv("HH_CLIENT_ID"),
		clientSecret: os.Getenv("HH_CLIENT_SECRET"),
		client:       resty.New(),
	}
}

// ExchangeCodeForToken обменивает `code` на `access_token`
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

// GetUserID получает ID текущего пользователя
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

// GetResume получает резюме по ID
func (c *HHClient) GetResume(accessToken, resumeID string) (*models.Resume, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(c.apiURL + fmt.Sprintf(constants.Resume, resumeID))

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить резюме: %s", resp.String())
	}

	var resume models.Resume
	if err := json.Unmarshal(resp.Body(), &resume); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &resume, nil
}

// GetResumes получает список резюме пользователя
func (c *HHClient) GetResumes(accessToken string) (*models.APIResumeResponse, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(c.apiURL + constants.ResumesMine)

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить список резюме: %s", resp.String())
	}

	var resumes models.APIResumeResponse
	if err := json.Unmarshal(resp.Body(), &resumes); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &resumes, nil
}

// GetVacancyByID получает вакансию по ID
func (c *HHClient) GetVacancyByID(accessToken, vacancyID string) (*models.Vacancy, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(c.apiURL + fmt.Sprintf(constants.Vacancy, vacancyID))

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить вакансию: %s", resp.String())
	}

	var vacancy models.Vacancy
	if err := json.Unmarshal(resp.Body(), &vacancy); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &vacancy, nil
}

func (c *HHClient) GetUserApplications(accessToken string) ([]models.ApplicationItem, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
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

func (c *HHClient) GetUserFirstFoundedApplication(accessToken string) (*models.APIApplicationsResponse, error) {
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
