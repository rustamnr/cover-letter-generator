package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

// HHService отвечает за работу с API hh.ru
type HHService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	apiURL       string
	tokenURL     string
	client       *resty.Client
}

// NewHHService создает новый экземпляр HHService
func NewHHService(clientID, clientSecret, redirectURI string) *HHService {
	return &HHService{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		apiURL:       "https://api.hh.ru",
		tokenURL:     "https://hh.ru/oauth/token",
		client:       resty.New(),
	}
}

// GetClient возвращает HTTP-клиент для запросов к API hh.ru
func (s *HHService) GetClient() *resty.Client {
	return s.client
}

// GetAPIURL возвращает базовый URL API hh.ru
func (s *HHService) GetAPIURL() string {
	return s.apiURL
}

// GetAuthURL возвращает ссылку для авторизации пользователя
func (s *HHService) GetAuthURL() string {
	return constants.GetAuthURL(s.clientID, s.redirectURI)
}

// ExchangeCodeForToken обменивает `code` на `access_token`
func (s *HHService) ExchangeCodeForToken(code string) (string, error) {
	resp, err := s.client.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     s.clientID,
			"client_secret": s.clientSecret,
			"code":          code,
			"redirect_uri":  s.redirectURI,
		}).
		Post(s.tokenURL)

	if err != nil {
		return "", errors.New("ошибка запроса к hh.ru")
	}

	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("не удалось получить токен: " + resp.String())
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &tokenData); err != nil {
		return "", err
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		return "", errors.New("access_token не найден")
	}

	return accessToken, nil
}

// GetUserID получает ID текущего пользователя
func (s *HHService) GetUserID(accessToken string) (string, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(s.apiURL + constants.Me)

	if err != nil || resp.StatusCode() != http.StatusOK {
		return "", errors.New("ошибка запроса к API hh.ru")
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &userData); err != nil {
		return "", err
	}

	userID, ok := userData["id"].(string)
	if !ok {
		return "", errors.New("не удалось получить user_id")
	}

	return userID, nil
}

func (s *HHService) GetResume(accessToken string, resumeID string) (*models.Resume, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(s.apiURL + fmt.Sprintf(constants.Resume, resumeID))
	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, err
	}
	var resume models.Resume
	if err := json.Unmarshal(resp.Body(), &resume); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &resume, nil
}

func (s *HHService) GetResumes(accessToken string) (*models.APIResumeResponse, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(s.apiURL + constants.ResumesMine)

	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, errors.New("error fetching resumes: " + resp.String())
	}

	var resumes models.APIResumeResponse
	if err := json.Unmarshal(resp.Body(), &resumes); err != nil {
		return nil, err
	}

	return &resumes, nil
}

func (s *HHService) GetVacancyByID(accessToken string, vacancyID string) (*models.Vacancy, error) {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(s.apiURL + fmt.Sprintf(constants.Vacancy, vacancyID))

	if err != nil || resp.StatusCode() != http.StatusOK {
		return nil, errors.New("error fetching vacancy: " + resp.String())
	}

	var vacancy models.Vacancy
	if err := json.Unmarshal(resp.Body(), &vacancy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return &vacancy, nil
}
