package models

import (
	"encoding/gob"
	"fmt"
	"strings"
)

func init() {
	gob.Register([]SessionResume{})
}

type SessionResume struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type APIResumeResponse struct {
	Items []Resume `json:"items"`
}

type Resume struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	MiddleName        *string            `json:"middle_name,omitempty"`
	Age               *int               `json:"age,omitempty"`
	Gender            *Gender            `json:"gender,omitempty"`
	Area              Area               `json:"area"`
	Salary            *Salary            `json:"salary,omitempty"`
	Contact           []Contact          `json:"contact"`
	Experience        []Experience       `json:"experience"`
	Education         EducationInfo      `json:"education"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	AlternateURL      string             `json:"alternate_url"`
	Status            Status             `json:"status"`
	TotalExperience   TotalExperience    `json:"total_experience"`
	ProfessionalRoles []ProfessionalRole `json:"professional_roles,omitempty"`
	Photo             *Photo             `json:"photo,omitempty"`
	KeySkills         []KeySkill         `json:"key_skills,omitempty"`
	Languages         []Language         `json:"languages,omitempty"`

	// Resume-specific fields
	Progress           *Progress       `json:"_progress,omitempty"`
	Access             *Access         `json:"access,omitempty"`
	Actions            *Actions        `json:"actions,omitempty"`
	BirthDate          *string         `json:"birth_date,omitempty"`
	Blocked            bool            `json:"blocked"`
	BusinessTrip       *BusinessTrip   `json:"business_trip_readiness,omitempty"`
	CanPublishOrUpdate bool            `json:"can_publish_or_update"`
	Certificate        []interface{}   `json:"certificate,omitempty"`
	Citizenship        []Area          `json:"citizenship,omitempty"`
	Creds              *Creds          `json:"creds,omitempty"`
	Download           *Download       `json:"download,omitempty"`
	DriverLicense      []DriverLicense `json:"driver_license_types,omitempty"`
	Employments        []Employment    `json:"employments,omitempty"`
	Finished           bool            `json:"finished"`
	HasVehicle         *bool           `json:"has_vehicle,omitempty"`
	Metro              *string         `json:"metro,omitempty"`
	ModerationNote     []interface{}   `json:"moderation_note,omitempty"`
	NextPublishAt      string          `json:"next_publish_at,omitempty"`
	Portfolio          []interface{}   `json:"portfolio,omitempty"`
	PublishURL         string          `json:"publish_url,omitempty"`
	Recommendation     []interface{}   `json:"recommendation,omitempty"`
	Relocation         *Relocation     `json:"relocation,omitempty"`
	ResumeLocale       *Locale         `json:"resume_locale,omitempty"`
	Schedules          []Schedule      `json:"schedules,omitempty"`
	Site               []interface{}   `json:"site,omitempty"`
	SkillSet           []interface{}   `json:"skill_set,omitempty"`
	Skills             string          `json:"skills,omitempty"`
	TravelTime         *TravelTime     `json:"travel_time,omitempty"`
	ViewsURL           string          `json:"views_url,omitempty"`
	WorkTicket         []interface{}   `json:"work_ticket,omitempty"`
}

type Gender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TotalExperience struct {
	Months int `json:"months"`
}

type Progress struct {
	Mandatory   []any `json:"mandatory"`
	Percentage  int   `json:"percentage"`
	Recommended []any `json:"recommended"`
}

type Recommended struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Access struct {
	Type AccessType `json:"type"`
}

type AccessType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Actions struct {
	Download Download `json:"download"`
}

type Download struct {
	PDF DownloadLink `json:"pdf"`
	RTF DownloadLink `json:"rtf"`
}

type DownloadLink struct {
	URL string `json:"url"`
}

type BusinessTrip struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Creds struct {
	Answers             map[string]CredAnswer   `json:"answers"`
	QuestionToAnswerMap map[string][]string     `json:"question_to_answer_map"`
	Questions           map[string]CredQuestion `json:"questions"`
}

type CredAnswer struct {
	AnswerGroup       string        `json:"answer_group"`
	AnswerID          string        `json:"answer_id"`
	AskQuestionsAfter []interface{} `json:"ask_questions_after"`
	Description       *string       `json:"description"`
	PositiveTitle     string        `json:"positive_title"`
	SkipAtResult      bool          `json:"skip_at_result"`
	Title             string        `json:"title"`
}

type CredQuestion struct {
	Description     *string  `json:"description"`
	IsActive        bool     `json:"is_active"`
	PossibleAnswers []string `json:"possible_answers"`
	QuestionID      string   `json:"question_id"`
	QuestionTitle   string   `json:"question_title"`
	QuestionType    string   `json:"question_type"`
	Required        bool     `json:"required"`
	SkipTitleAtView bool     `json:"skip_title_at_view"`
	ViewTitle       *string  `json:"view_title"`
}

type Relocation struct {
	Area     []interface{}  `json:"area"`
	District []interface{}  `json:"district"`
	Type     RelocationType `json:"type"`
}

type RelocationType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Locale struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TravelTime struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DriverLicense struct {
	ID string `json:"id"`
}

func (r *Resume) ConvertToText() string {
	var sb strings.Builder
	sb.WriteString(r.getBasicInfo())
	sb.WriteString(r.getLocationInfo())
	sb.WriteString(r.getContactInfo())
	sb.WriteString(r.getWorkExperienceInfo())
	sb.WriteString(r.getEducationInfo())
	sb.WriteString(r.getSalaryInfo())
	return sb.String()
}

func (r *Resume) getBasicInfo() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Резюме: %s\n", r.Title))
	sb.WriteString(fmt.Sprintf("Имя: %s %s\n", r.FirstName, r.LastName))
	if r.MiddleName != nil {
		sb.WriteString(fmt.Sprintf("Отчество: %s\n", *r.MiddleName))
	}
	if r.Age != nil {
		sb.WriteString(fmt.Sprintf("Возраст: %d\n", *r.Age))
	}
	if r.Gender != nil {
		sb.WriteString(fmt.Sprintf("Пол: %s\n", r.Gender.Name))
	}
	sb.WriteString(fmt.Sprintf("Дата создания: %s\n", r.CreatedAt))
	sb.WriteString(fmt.Sprintf("Дата обновления: %s\n", r.UpdatedAt))
	sb.WriteString(fmt.Sprintf("Статус: %s\n", r.Status.Name))
	return sb.String()
}

func (r *Resume) getLocationInfo() string {
	return fmt.Sprintf("Местоположение: %s\n", r.Area.Name)
}

func (r *Resume) getContactInfo() string {
	var sb strings.Builder
	sb.WriteString("Контакты:\n")
	for _, contact := range r.Contact {
		sb.WriteString(fmt.Sprintf("%s\n", r.getContactText(contact)))
	}
	return sb.String()
}

func (r *Resume) getContactText(contact Contact) string {
	var contactInfo string

	switch v := contact.Value.(type) {
	case string:
		contactInfo = fmt.Sprintf("%s: %s", contact.Type.Name, v)
	case map[string]interface{}:
		phone := PhoneValue{}
		if city, ok := v["city"].(string); ok {
			phone.City = city
		}
		if country, ok := v["country"].(string); ok {
			phone.Country = country
		}
		if number, ok := v["number"].(string); ok {
			phone.Number = number
		}
		if formatted, ok := v["formatted"].(string); ok {
			phone.Formatted = formatted
		}
		contactInfo = fmt.Sprintf("Телефон: %s", phone.Formatted)
	}

	if contact.Comment != nil {
		contactInfo += fmt.Sprintf(" (%s)", *contact.Comment)
	}

	return contactInfo
}

func (r *Resume) getWorkExperienceInfo() string {
	var sb strings.Builder
	sb.WriteString("Опыт работы:\n")
	for _, exp := range r.Experience {
		endDate := "по настоящее время"
		if exp.EndDate != nil {
			endDate = *exp.EndDate
		}
		sb.WriteString(fmt.Sprintf("%s - %s (%s - %s)\n", exp.Position, exp.Company, exp.StartDate, endDate))
		if exp.Description != "" {
			sb.WriteString(fmt.Sprintf("Описание: %s\n", exp.Description))
		}
	}
	sb.WriteString(fmt.Sprintf("Общий опыт: %d месяцев\n", r.TotalExperience.Months))
	return sb.String()
}

func (r *Resume) getEducationInfo() string {
	var sb strings.Builder
	sb.WriteString("Образование:\n")
	sb.WriteString(fmt.Sprintf("Уровень: %s\n", r.Education.Level.Name))
	for _, edu := range r.Education.Primary {
		sb.WriteString(fmt.Sprintf("%s", edu.Name))
		if edu.Organization != "" {
			sb.WriteString(fmt.Sprintf(", %s", edu.Organization))
		}
		if edu.Year != 0 {
			sb.WriteString(fmt.Sprintf(", %d год", edu.Year))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r *Resume) getSalaryInfo() string {
	if r.Salary != nil && r.Salary.Amount != nil {
		return fmt.Sprintf("Зарплата: %d %s\n", *r.Salary.Amount, r.Salary.Currency)
	}
	return "Зарплата: Не указана\n"
}
