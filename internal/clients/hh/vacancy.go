package clients_hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rustamnr/cover-letter-generator/internal/constants"
	"github.com/rustamnr/cover-letter-generator/internal/logger"
	"github.com/rustamnr/cover-letter-generator/internal/models"
)

func (c *HHClient) GetVacancyByID(vacancyID string) (*models.Vacancy, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		Get(c.apiURL + fmt.Sprintf(constants.Vacancy, vacancyID))
	if err != nil {
		return nil, fmt.Errorf("failed to get vacancy: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful response from hh.ru: %s", resp.String())
	}

	var vacancy models.Vacancy
	if err := json.Unmarshal(resp.Body(), &vacancy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	vacancy.Description = cleanHTML(vacancy.Description)
	if vacancy.BrandedDescription != nil {
		*vacancy.BrandedDescription = cleanHTML(*vacancy.BrandedDescription)
	}

	logger.Debugf("vacancy: %+v", vacancy)
	return &vacancy, nil
}

// GetSimilarVacancies retrieves a list of similar vacancies for a given resume ID
func (c *HHClient) GetSimilarVacancies(
	resumeID string, queryParams map[string]string) ([]models.Vacancy, error) {
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.accessToken).
		SetQueryParams(queryParams).
		Get(c.apiURL + fmt.Sprintf("/resumes/%s/similar_vacancies", resumeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get similar vacancies: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful response from hh.ru: %s", resp.String())
	}

	var similarVacancies models.VacanciesResponse[models.Vacancy]
	if err := json.Unmarshal(resp.Body(), &similarVacancies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return similarVacancies.Items, nil
}

func (c *HHClient) GetFirstSuitableVacancy(resumeID string) (*models.Vacancy, error) {
	firstSimilarVacancy, err := c.GetSimilarVacancies(resumeID, map[string]string{"per_page": "1"})
	if err != nil {
		return nil, fmt.Errorf("failed to get similar vacancies: %w", err)
	}
	if len(firstSimilarVacancy) != 1 {
		return nil, errors.New("vacancies count is not 1")
	}

	return &firstSimilarVacancy[0], nil
}
