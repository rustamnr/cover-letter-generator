package clients_hh

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

func (c *HHClient) GetResume(resumeID string) (*models.Resume, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
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

func (c *HHClient) GetResumes() (*models.ResumesResponse, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		Get(c.apiURL + constants.ResumesMine)

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API hh.ru: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить список резюме: %s", resp.String())
	}

	var resumes models.ResumesResponse
	if err := json.Unmarshal(resp.Body(), &resumes); err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %w", err)
	}

	return &resumes, nil
}
