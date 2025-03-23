package models

// APIApplicationsResponse представляет список откликов пользователя
type APIApplicationsResponse struct {
	Items []ApplicationItem `json:"items"`
	Page  int               `json:"page"`
	Pages int               `json:"pages"`
}

// ApplicationItem представляет отклик на вакансию
type ApplicationItem struct {
	ID        string   `json:"id"`
	Vacancy   Vacancy  `json:"vacancy"`
	Employer  Employer `json:"employer"`
	Status    Status   `json:"status"`
	CreatedAt string   `json:"created_at"`
}

// Vacancy содержит информацию о вакансии
type Vacancy struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	AlternateURL string `json:"alternate_url"`
}

// Employer содержит информацию о работодателе
type Employer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Status описывает статус отклика
type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
