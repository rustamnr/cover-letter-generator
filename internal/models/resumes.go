package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
	Salary     *ResumeSalary    `json:"salary"`
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
	Company    string     `json:"company"`
	Position   string     `json:"position"`
	StartDate  string     `json:"start"`
	EndDate    *string    `json:"end"`
	Area       *Location  `json:"area,omitempty"`
	Industry   *Industry  `json:"industry,omitempty"`
	Industries []Industry `json:"industries,omitempty"`
}

type Industry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EducationInfo struct {
	Level   EducationLevel `json:"level"`
	Primary []Education    `json:"primary"`
}

type EducationLevel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Education struct {
	Name string `json:"name"`
}

type ResumeSalary struct {
	Amount   *int   `json:"amount"`
	Currency string `json:"currency"`
}

func (r *HHResume) ConvertToText() string {
	var sb strings.Builder
	sb.WriteString(r.getBasicInfo())
	sb.WriteString(r.getLocationInfo())
	sb.WriteString(r.getContactInfo())
	sb.WriteString(r.getWorkExperienceInfo())
	sb.WriteString(r.getEducationInfo())
	sb.WriteString(r.getSalaryInfo())
	return sb.String()
}

func (r *HHResume) getBasicInfo() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Резюме:\n"))
	sb.WriteString(fmt.Sprintf("Имя: %s %s\n", r.FirstName, r.LastName))
	if r.MiddleName != nil {
		sb.WriteString(fmt.Sprintf("Отчество: %s\n", *r.MiddleName))
	}
	sb.WriteString(fmt.Sprintf("Должность: %s\n", r.Title))
	sb.WriteString(fmt.Sprintf("Дата создания: %s\n", r.CreatedAt))
	sb.WriteString(fmt.Sprintf("Дата обновления: %s\n", r.UpdatedAt))
	return sb.String()
}

func (r *HHResume) getLocationInfo() string {
	return fmt.Sprintf("Местоположение: %s\n", r.Area.Name)
}

func (r *HHResume) getContactInfo() string {
	var sb strings.Builder
	sb.WriteString("Контакты:\n")
	for _, contact := range r.Contact {
		sb.WriteString(fmt.Sprintf("%s\n", r.getContactText(contact)))
	}
	return sb.String()
}

func (r *HHResume) getContactText(contact Contact) string {
	var contactInfo string
	if contact.Type.Name == "phone" {
		var phone PhoneValue
		if err := json.Unmarshal(contact.Value, &phone); err == nil {
			contactInfo = fmt.Sprintf("Телефон: +%s (%s) %s", phone.Country, phone.City, phone.Number)
		}
	} else if contact.Type.Name == "email" {
		var email string
		if err := json.Unmarshal(contact.Value, &email); err == nil {
			contactInfo = fmt.Sprintf("Email: %s", email)
		}
	}
	return contactInfo
}

func (r *HHResume) getWorkExperienceInfo() string {
	var sb strings.Builder
	sb.WriteString("Опыт работы:\n")
	for _, exp := range r.Experience {
		sb.WriteString(fmt.Sprintf("%s - %s (%s - %s)\n", exp.Position, exp.Company, exp.StartDate, *exp.EndDate))
	}
	return sb.String()
}

func (r *HHResume) getEducationInfo() string {
	var sb strings.Builder
	sb.WriteString("Образование:\n")
	for _, edu := range r.Education.Primary {
		sb.WriteString(fmt.Sprintf("%s\n", edu.Name))
	}
	return sb.String()
}

func (r *HHResume) getSalaryInfo() string {
	if r.Salary != nil {
		return fmt.Sprintf("Зарплата: %d %s\n", *r.Salary.Amount, r.Salary.Currency)
	}
	return "Зарплата: Не указана\n"
}
