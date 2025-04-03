package models

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Contact struct {
	Comment          *string       `json:"comment,omitempty"`
	NeedVerification bool          `json:"need_verification,omitempty"`
	Preferred       bool           `json:"preferred"`
	Type            ContactType    `json:"type"`
	Value           interface{}    `json:"value"` // string or PhoneValue
	Verified        bool           `json:"verified,omitempty"`
}

type ContactType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PhoneValue struct {
	City      string `json:"city,omitempty"`
	Country   string `json:"country,omitempty"`
	Formatted string `json:"formatted,omitempty"`
	Number    string `json:"number,omitempty"`
}

type Salary struct {
	From     *int   `json:"from,omitempty"`
	To       *int   `json:"to,omitempty"`
	Amount   *int   `json:"amount,omitempty"` // For resumes
	Currency string `json:"currency"`
	Gross    bool   `json:"gross,omitempty"`
}

type Employment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Experience struct {
	ID          string     `json:"id,omitempty"` // For vacancies
	Name        string     `json:"name,omitempty"`
	Company     string     `json:"company,omitempty"` // For resumes
	Position    string     `json:"position,omitempty"`
	StartDate   string     `json:"start,omitempty"`
	EndDate     *string    `json:"end,omitempty"`
	Description string     `json:"description,omitempty"`
	Area        *Area      `json:"area,omitempty"`
	Industry    *Industry  `json:"industry,omitempty"`
	Industries  []Industry `json:"industries,omitempty"`
}

type Industry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Education struct {
	ID                string  `json:"id,omitempty"`
	Name              string  `json:"name"`
	NameID            *string `json:"name_id,omitempty"`
	Organization      string  `json:"organization,omitempty"`
	OrganizationID    *string `json:"organization_id,omitempty"`
	Result            string  `json:"result,omitempty"`
	ResultID          *string `json:"result_id,omitempty"`
	UniversityAcronym string  `json:"university_acronym,omitempty"`
	Year              int     `json:"year,omitempty"`
}

type EducationLevel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EducationInfo struct {
	Level       EducationLevel `json:"level"`
	Primary     []Education    `json:"primary"`
	Additional  []interface{}  `json:"additional,omitempty"`
	Attestation []interface{}  `json:"attestation,omitempty"`
	Elementary  []interface{}  `json:"elementary,omitempty"`
}

type Language struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level Level  `json:"level,omitempty"`
}

type Level struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProfessionalRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Schedule struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Photo struct {
	Small    string `json:"90,omitempty"`
	Medium   string `json:"240,omitempty"`
	Original string `json:"original,omitempty"`
	Size40   string `json:"40,omitempty"`
	Size100  string `json:"100,omitempty"`
	Size500  string `json:"500,omitempty"`
	ID       string `json:"id,omitempty"`
}

type KeySkill struct {
	Name string `json:"name"`
}