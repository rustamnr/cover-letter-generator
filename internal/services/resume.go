package services

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/models"
)

// ResumeService содержит методы для работы с резюме
type ResumeService struct {
	APIClient *http.Client
	APIUrl    string
}

// NewResumeService создает новый сервис
func NewResumeService(apiURL string) *ResumeService {
	return &ResumeService{
		APIClient: &http.Client{},
		APIUrl:    apiURL,
	}
}

// GetResume получает резюме пользователя с hh.ru
func (s *ResumeService) GetResume(accessToken string) (*models.Resume, error) {
	req, err := http.NewRequest("GET", s.APIUrl+"/resumes/mine", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.APIClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch resume")
	}

	var resumeData models.Resume
	if err := json.NewDecoder(resp.Body).Decode(&resumeData); err != nil {
		return nil, err
	}

	return &resumeData, nil
}
