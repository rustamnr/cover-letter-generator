package models

import "encoding/json"

type APIResumeResponse struct {
	Items []HHResume `json:"items"`
}

type HHResume struct {
	ID         string           `json:"id"`
	FirstName  string           `json:"first_name"`
	LastName   string           `json:"last_name"`
	MiddleName *string          `json:"middle_name"`
	Title      string           `json:"title"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Area       Location         `json:"area"`
	Contact    []Contact        `json:"contact"`
	Experience []WorkExperience `json:"experience"`
	Education  EducationInfo    `json:"education"`
	Salary     *Salary          `json:"salary"`
}

type Location struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Contact struct {
	Type        ContactType     `json:"type"`
	Value       json.RawMessage `json:"value"` // Может быть строкой (email) или объектом (телефон)
	Preferred   bool            `json:"preferred"`
	Verified    bool            `json:"verified"`
	Comment     *string         `json:"comment"`
	ParsedValue string          `json:"-"` // Распарсенное значение (email или телефон)
}

type ContactType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Структура для хранения телефонных данных
type PhoneValue struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Number  string `json:"number"`
}

type WorkExperience struct {
	Company    string      `json:"company"`
	Position   string      `json:"position"`
	StartDate  string      `json:"start"`
	EndDate    *string     `json:"end"`
	Area       *Location   `json:"area,omitempty"`
	Industry   *Industry   `json:"industry,omitempty"`
	Industries []Industry  `json:"industries,omitempty"`
}

type Industry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EducationInfo struct {
	Level   EducationLevel `json:"level"`
	Primary []Education   `json:"primary"`
}

type EducationLevel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Education struct {
	Name string `json:"name"`
}

type Salary struct {
	Amount   *int   `json:"amount"`
	Currency string `json:"currency"`
}

